package webcrawler

import (
	"fmt"
	"strings"
	"time"

	"golang.org/x/net/html"
)

//Item represents extracted item
type Item struct {
	ItemName    string
	URL         *URL
	TimeQueried string
	DateQueried string
	ItemDetails map[string]string
}

//ScrapeItemConfig configuration used to extract item from html token
type ScrapeItemConfig struct {
	ItemName    string                            `json:"ItemName"`
	ItemToGet   ExtractFromTokenConfig            `json:"ItemToGet"`
	ItemDetails map[string]ExtractFromTokenConfig `json:"ItemDetails"`
}

// ExtractItemWithScrapItemConfig extracts item from html token using a list of scrape item config which allows
// for selective extraction.
func ExtractItemWithScrapItemConfig(t html.Token, z *html.Tokenizer, itemTagsToCheck map[string]bool, scrapeItemConfig []ScrapeItemConfig) (Item, error) {
	if itemTagsToCheck[t.Data] {
		for _, scrapeItemConfig := range scrapeItemConfig {
			HTTPAttributeValueFromToken, err := extractAttributeValue(t, scrapeItemConfig.ItemToGet.Attribute)
			if err != nil {
				return Item{}, err
			}
			if (t.Data == scrapeItemConfig.ItemToGet.Tag && HTTPAttributeValueFromToken == scrapeItemConfig.ItemToGet.AttributeValue) || (scrapeItemConfig.ItemToGet.Attribute == "" && scrapeItemConfig.ItemToGet.AttributeValue == "") {
				extractedItem, err := parseTokenForItemDetails(t, z, scrapeItemConfig)
				if err != nil {
					return Item{}, err
				}
				return extractedItem, nil
			}
		}
	}
	return Item{}, fmt.Errorf("unable to extract item with scrape item config")
}

// parseTokenForItemDetails parses html block and extracts item details using the scrape item config from the tags, attributes,
// and comments of the html elemements. Given the start token of an html block, parse that block by pushing each start token onto the stack,
// when encountering end tokens, pop tokens off of the stack. Since start tokens have corresponding end tokens, the end of the block is reached
// when the length of the stack is 0. Returns the item with the item details.
func parseTokenForItemDetails(token html.Token, z *html.Tokenizer, scrapeItemConfig ScrapeItemConfig) (Item, error) {
	var (
		tokenType             html.TokenType
		currentToken          html.Token
		tagStack              stack
		itemDetailTagsToCheck map[string]bool
		item                  Item = Item{
			ItemName:    scrapeItemConfig.ItemName,
			ItemDetails: make(map[string]string),
			DateQueried: strings.Split(time.Now().String(), " ")[0],
			TimeQueried: strings.Split(time.Now().String(), " ")[1],
		}
	)

	if token.Type != html.StartTagToken {
		return item, fmt.Errorf("unable to parse item, not a start tag")
	}
	itemDetailTagsToCheck, _ = generateItemDetailsTagsToCheckMap(itemDetailTagsToCheck, scrapeItemConfig)
	if err != nil {
		return item, nil
	}
	tagStack.push(token)
	for len(tagStack) != 0 {
		tokenType = z.Next()
		currentToken = z.Token()
		switch {
		case tokenType == html.StartTagToken:
			tagStack.push(currentToken)
			if itemDetailTagsToCheck[currentToken.Data] {
				for itemDetailName, itemDetails := range scrapeItemConfig.ItemDetails {
					HTTPAttributeValueFromToken, _ := extractAttributeValue(currentToken, itemDetails.Attribute)
					if (itemDetails.Tag == currentToken.Data && itemDetails.AttributeValue == HTTPAttributeValueFromToken) || (itemDetails.Tag == currentToken.Data && itemDetails.Attribute == "" && itemDetails.AttributeValue == "") {
						if _, exist := item.ItemDetails[itemDetailName]; exist {
							return item, nil
						}
						if itemDetails.AttributeToGet != "" {
							HTTPAttributeValueFromToken, _ = extractAttributeValue(currentToken, itemDetails.AttributeToGet)
							if !IsEmpty(itemDetails.FormatAttributeConfiguration) {
								HTTPAttributeValueFromToken = formatURL(HTTPAttributeValueFromToken, itemDetails.FormatAttributeConfiguration)
							}
							item.ItemDetails[itemDetailName] = HTTPAttributeValueFromToken
						} else {
							if itemDetails.SkipToken != 0 {
								for itemDetails.SkipToken >= 0 {
									tokenType = z.Next()
									if tokenType == html.TextToken {
										itemDetails.SkipToken--
										continue
									}
								}
							} else {
								for tokenType != html.TextToken {
									tokenType = z.Next()
								}
							}
							currentToken = z.Token()
							str := currentToken.String()
							if !IsEmpty(itemDetails.FormatAttributeConfiguration) {
								HTTPAttributeValueFromToken = formatURL(str, itemDetails.FormatAttributeConfiguration)
							}

							if !IsEmpty(itemDetails.ItemFilterConfiguration) {
								if !Validate(str, &itemDetails.ItemFilterConfiguration) {
									return Item{}, nil
								}
							}
							item.ItemDetails[itemDetailName] = str
						}

					}
				}
			}
		case tokenType == html.EndTagToken:
			tagStack.pop()
		case tokenType == html.ErrorToken:
			tagStack.pop()
			return item, nil
		default:
		}
	}
	return item, nil
}

// generateItemDetailsTagsToCheckMap generates a map of tags to check given the scrape item configuration. The map is used to check
// whether or not the html element should be used to extract from
func generateItemDetailsTagsToCheckMap(itemDetailTagsToCheck map[string]bool, scrapeItemConfig ScrapeItemConfig) (map[string]bool, error) {
	for _, item := range scrapeItemConfig.ItemDetails {
		if len(itemDetailTagsToCheck) == 0 {
			itemDetailTagsToCheck = make(map[string]bool)
		}
		itemDetailTagsToCheck[item.Tag] = true
	}
	return itemDetailTagsToCheck, nil
}
