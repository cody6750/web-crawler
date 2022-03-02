package webcrawler

import (
	"fmt"

	"golang.org/x/net/html"
)

//ExtractFromTokenConfig ...
type ExtractFromTokenConfig struct {
	ItemFilterConfiguration      FilterConfiguration `json:"FilterConfiguration"`
	FormatAttributeConfiguration FormatURLConfig     `json:"FormatAttributeConfiguration"`
	SkipToken                    int                 `json:"SkipToken"`
	Tag                          string              `json:"Tag"`
	Attribute                    string              `json:"Attribute"`
	AttributeValue               string              `json:"AttributeValue"`
	AttributeToGet               string              `json:"AttributeToGet"`
}

func extractAttributeValue(token html.Token, attributeToGet string) (attributeValue string, err error) {
	for _, a := range token.Attr {
		if a.Key == attributeToGet {
			attributeValue = a.Val
			return attributeValue, nil
		}
	}
	if attributeValue == "" {
		return attributeValue, fmt.Errorf("unable to get attribute value from %v", attributeToGet)
	}
	return attributeValue, nil
}
