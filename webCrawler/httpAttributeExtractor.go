package webcrawler

import (
	"errors"

	"golang.org/x/net/html"
)

var (
	errEmptyParameter error = errors.New("Empty parameter")
)

//ExtractFromHTMLConfiguration ...
type ExtractFromHTMLConfiguration struct {
	TagToCheck            string
	AttributeToCheck      string
	AttributeValueToCheck string
}

func getHTTPAttributeValueFromToken(token html.Token, attributeToGet string) (attributeValue string, err error) {
	if attributeToGet == "" {
		return "empty string", errEmptyParameter
	}
	for _, a := range token.Attr {
		if a.Key == attributeToGet {
			attributeValue = a.Val
			return attributeValue, nil
		}
	}
	if attributeValue == "" {
		return "empty string", errEmptyParameter
	}
	return "does not exist", nil
}
