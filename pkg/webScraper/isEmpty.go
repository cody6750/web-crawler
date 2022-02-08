package webcrawler

import (
	"reflect"

	"golang.org/x/net/html"
)

//IsEmpty ...
func IsEmpty(i interface{}) bool {
	switch i := i.(type) {
	case ExtractFromTokenConfig:
		return reflect.DeepEqual(i, ExtractFromTokenConfig{})
	case html.Token:
		return i.Data == "" && len(i.Attr) == 0
	case FormatURLConfig:
		return reflect.DeepEqual(i, FormatURLConfig{})
	case []ScrapeItemConfig:
		return len(i) == 0
	case []ScrapeURLConfig:
		return len(i) == 0
	case map[string]ExtractFromTokenConfig:
		return len(i) == 0
	case map[string]struct{}:
		return len(i) == 0
	default:
		return false
	}
}
