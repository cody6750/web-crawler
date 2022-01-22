package webcrawler

import (
	"golang.org/x/net/html"
)

//ExtractFromHTMLConfiguration ...
type ExtractFromHTMLConfiguration struct {
	Tag            string
	Attribute      string
	AttributeValue string
	AttributeToGet string
}

func getHTTPAttributeValueFromToken(token html.Token, attributeToGet string) (attributeValue string, err error) {
	if attributeToGet == "" {
		return attributeValue, errEmptyParameter
	}
	for _, a := range token.Attr {
		if a.Key == attributeToGet {
			attributeValue = a.Val
			return attributeValue, nil
		}
	}
	if attributeValue == "" {
		return attributeValue, errEmptyParameter
	}
	return attributeValue, nil
}
