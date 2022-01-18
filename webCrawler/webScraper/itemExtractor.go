package webcrawler

import (
	"errors"
	"strings"
	"time"

	"golang.org/x/net/html"
)

//ScrapeItemConfiguration ...
type ScrapeItemConfiguration struct {
	ItemName    string
	ItemToGet   ExtractFromHTMLConfiguration
	ItemDetails map[string]ExtractFromHTMLConfiguration
}

//ExtractItemWithScrapItemConfiguration ...
func ExtractItemWithScrapItemConfiguration(t html.Token, z *html.Tokenizer, itemTagsToCheck map[string]bool, scrapeItemConfiguration []ScrapeItemConfiguration) (Item, error) {
	if itemTagsToCheck[t.Data] {
		for _, scrapeItemConfiguration := range scrapeItemConfiguration {
			HTTPAttributeValueFromToken, _ := getHTTPAttributeValueFromToken(t, scrapeItemConfiguration.ItemToGet.Attribute)
			if (t.Data == scrapeItemConfiguration.ItemToGet.Tag && HTTPAttributeValueFromToken == scrapeItemConfiguration.ItemToGet.AttributeValue) || (scrapeItemConfiguration.ItemToGet.Attribute == "" && scrapeItemConfiguration.ItemToGet.AttributeValue == "") {
				extractedItem, _ := parseTokenForItemDetails(t, z, scrapeItemConfiguration)
				return extractedItem, nil
			}
		}
	}
	return Item{}, errors.New("")
}

func parseTokenForItemDetails(token html.Token, z *html.Tokenizer, scrapeItemConfiguration ScrapeItemConfiguration) (Item, error) {
	var (
		tokenType             html.TokenType
		currentToken          html.Token
		tagStack              stack
		itemDetailTagsToCheck map[string]bool
		item                  Item = Item{
			ItemName:    scrapeItemConfiguration.ItemName,
			ItemDetails: make(map[string]string),
			DateQueried: strings.Split(time.Now().String(), " ")[0],
			TimeQueried: strings.Split(time.Now().String(), " ")[1],
		}
	)

	if token.Type != html.StartTagToken {
		return item, errors.New("Unable to parse item, not a start tag")
	}
	itemDetailTagsToCheck, _ = generateItemDetailsTagsToCheckMap(itemDetailTagsToCheck, scrapeItemConfiguration)
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
				for itemDetailName, itemDetails := range scrapeItemConfiguration.ItemDetails {
					HTTPAttributeValueFromToken, _ := getHTTPAttributeValueFromToken(currentToken, itemDetails.Attribute)
					if (itemDetails.Tag == currentToken.Data && itemDetails.AttributeValue == HTTPAttributeValueFromToken) || (itemDetails.Tag == currentToken.Data && itemDetails.Attribute == "" && itemDetails.AttributeValue == "") {
						if itemDetails.AttributeToGet != "" {
							HTTPAttributeValueFromToken, _ = getHTTPAttributeValueFromToken(currentToken, itemDetails.AttributeToGet)
							item.ItemDetails[itemDetailName] = HTTPAttributeValueFromToken
							continue
						} else {
							for tokenType != html.TextToken {
								tokenType = z.Next()
							}
							currentToken = z.Token()
							str := currentToken.String()
							item.ItemDetails[itemDetailName] = str
							continue
						}
					}
				}
			}
		case tokenType == html.EndTagToken:
			tagStack.pop()
		case tokenType == html.ErrorToken:
			tagStack.pop()
			// item.printJSON()
			return item, nil
		}
	}
	// item.printJSON()
	return item, nil
}

func generateItemDetailsTagsToCheckMap(itemDetailTagsToCheck map[string]bool, scrapeItemConfiguration ScrapeItemConfiguration) (map[string]bool, error) {
	if IsEmpty(scrapeItemConfiguration.ItemDetails) {
		return nil, errors.New("Item is empty")
	}
	for _, item := range scrapeItemConfiguration.ItemDetails {
		if len(itemDetailTagsToCheck) == 0 {
			itemDetailTagsToCheck = make(map[string]bool)
		}
		itemDetailTagsToCheck[item.Tag] = true
	}

	return itemDetailTagsToCheck, nil
}
