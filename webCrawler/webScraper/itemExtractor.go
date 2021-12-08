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
	URL         string
	ItemToGet   ExtractFromHTMLConfiguration
	ItemDetails map[string]ExtractFromHTMLConfiguration
}

//ExtractItemWithScrapItemConfiguration ...
func ExtractItemWithScrapItemConfiguration(t html.Token, z *html.Tokenizer, url string, itemTagsToCheck map[string]bool, scrapeItemConfiguration []ScrapeItemConfiguration) error {
	if itemTagsToCheck[t.Data] {
		for _, scrapeItemConfiguration := range scrapeItemConfiguration {
			HTTPAttributeValueFromToken, _ := getHTTPAttributeValueFromToken(t, scrapeItemConfiguration.ItemToGet.Attribute)
			if t.Data == scrapeItemConfiguration.ItemToGet.Tag && HTTPAttributeValueFromToken == scrapeItemConfiguration.ItemToGet.AttributeValue {
				scrapeItemConfiguration.URL = url
				parseTokenForItemDetails(t, z, scrapeItemConfiguration)
				return nil
			}
		}
	} else {
		return nil
	}
	return nil
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
			URL:         scrapeItemConfiguration.URL,
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
						} else {
							for tokenType != html.TextToken {
								tokenType = z.Next()
							}
							currentToken = z.Token()
							str := currentToken.String()
							item.ItemDetails[itemDetailName] = str
						}
					}
				}
			}
		case tokenType == html.EndTagToken:
			tagStack.pop()
		case tokenType == html.ErrorToken:
			item.printJSON()
			return item, nil
		}
	}
	item.printJSON()
	return item, nil
}
func isEmptyItemToGet(itemToget map[string]ExtractFromHTMLConfiguration) bool {
	return len(itemToget) == 0
}

func generateItemDetailsTagsToCheckMap(itemDetailTagsToCheck map[string]bool, scrapeItemConfiguration ScrapeItemConfiguration) (map[string]bool, error) {
	if isEmptyItemToGet(scrapeItemConfiguration.ItemDetails) {
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
