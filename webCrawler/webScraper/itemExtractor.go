package webcrawler

import (
	"errors"

	"golang.org/x/net/html"
)

//ScrapeItemConfiguration ...
type ScrapeItemConfiguration struct {
	ItemName  string
	URL       string
	ItemToget map[string]ExtractFromHTMLConfiguration
}

//Item ...
type Item struct {
	ItemName    string
	URL         string
	ItemDetails map[string]string
}

//ExtractItemWithScrapItemConfiguration ...
func ExtractItemWithScrapItemConfiguration(token html.Token, url string, itemTagsToCheck map[string]bool, scrapeItemConfiguration []ScrapeItemConfiguration) (Item, error) {
	item := Item{}
	if isEmptyToken(token) {
		return item, errors.New("Empty token")
	}
	// if isEmptyItemToGet(scrapeItemConfiguration.ItemToget) {
	// 	return item, errors.New("Empty ItemToGet")

	// }
	if itemTagsToCheck[token.Data] {
		for _, scrapeItemConfiguration := range scrapeItemConfiguration {
			for _, ItemToGet := range scrapeItemConfiguration.ItemToget {
				HTTPAttributeValueFromToken, _ := getHTTPAttributeValueFromToken(token, ItemToGet.AttributeToCheck)
				if ItemToGet.AttributeValueToCheck == HTTPAttributeValueFromToken {

				}
			}
		}
		// getHTTPAttributeValueFromToken(token, item.ItemToget)
	}
	return item, nil
}

func isEmptyItemToGet(itemToget map[string]ExtractFromHTMLConfiguration) bool {
	return len(itemToget) == 0
}
func (i *Item) printItemExtracted() {}
func (i *Item) getItem() string     { return i.ItemName }
func (i *Item) getURL() string      { return i.URL }
