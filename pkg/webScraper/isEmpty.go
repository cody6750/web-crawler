package webcrawler

import (
	"reflect"

	"golang.org/x/net/html"
)

//IsEmpty ...
func IsEmpty(i interface{}) bool {
	switch i := i.(type) {
	case ExtractFromHTMLConfiguration:
		return reflect.DeepEqual(i, ExtractFromHTMLConfiguration{})
	case html.Token:
		return i.Data == "" && len(i.Attr) == 0
	case FormatURLConfiguration:
		return reflect.DeepEqual(i, FormatURLConfiguration{})
	case []ScrapeItemConfig:
		return len(i) == 0
	case []ScrapeURLConfig:
		return len(i) == 0
	case map[string]ExtractFromHTMLConfiguration:
		return len(i) == 0
	case map[string]struct{}:
		return len(i) == 0
	default:
		return false
	}
}
