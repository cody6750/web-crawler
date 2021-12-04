package webcrawler

import (
	"encoding/json"
	"errors"
	"log"

	"golang.org/x/net/html"
)

type stack []html.Token

//ScrapeItemConfiguration ...
type ScrapeItemConfiguration struct {
	ItemName    string
	URL         string
	ItemToGet   ExtractFromHTMLConfiguration
	ItemDetails map[string]ExtractFromHTMLConfiguration
}

//Item ...
type Item struct {
	ItemName    string
	URL         string
	ItemDetails map[string]string
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
		item                  Item
	)
	if token.Type != html.StartTagToken {
		return item, errors.New("Unable to parse item, not a start tag")
	}
	itemDetailTagsToCheck, _ = generateItemDetailsTagsToCheckMap(itemDetailTagsToCheck, scrapeItemConfiguration)
	if err != nil {
		return item, nil
	}
	item.ItemName = scrapeItemConfiguration.ItemName
	item.URL = scrapeItemConfiguration.URL
	item.ItemDetails = make(map[string]string)
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
					if itemDetails.Tag == currentToken.Data && itemDetails.AttributeValue == HTTPAttributeValueFromToken {
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
		case tokenType == html.EndTagToken:
			tagStack.pop()
		case tokenType == html.ErrorToken:
			json, _ := json.MarshalIndent(item, "", "    ")
			log.Print("\n" + string(json))
			return item, nil
		}
	}
	json, _ := json.MarshalIndent(item, "", "    ")
	log.Print("\n" + string(json))
	return item, nil
}
func isEmptyItemToGet(itemToget map[string]ExtractFromHTMLConfiguration) bool {
	return len(itemToget) == 0
}
func (i *Item) printItemExtracted() {}
func (i *Item) getItem() string     { return i.ItemName }
func (i *Item) getURL() string      { return i.URL }

func (s *stack) push(t html.Token) {
	*s = append(*s, t)
}

func (s *stack) pop() (html.Token, bool) {
	if len(*s) == 0 {
		return html.Token{}, false
	}
	last := len(*s) - 1
	popped := (*s)[last]
	*s = (*s)[:last]
	return popped, true
}

func (s *stack) print() {
	log.Print(*s)
}

func generateItemDetailsTagsToCheckMap(itemDetailTagsToCheck map[string]bool, scrapeItemConfiguration ScrapeItemConfiguration) (map[string]bool, error) {
	if isEmptyItemToGet(scrapeItemConfiguration.ItemDetails) {
		return nil, errors.New("Item is empty")
	}
	// If item parameter is provided, create a map filled with tags that will be used to determine if processing is needed. To increase performance
	for _, item := range scrapeItemConfiguration.ItemDetails {
		if len(itemDetailTagsToCheck) == 0 {
			itemDetailTagsToCheck = make(map[string]bool)
		}
		itemDetailTagsToCheck[item.Tag] = true
	}

	return itemDetailTagsToCheck, nil
}
