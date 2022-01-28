package webcrawler

import "reflect"

//IsEmpty ...
func IsEmpty(i interface{}) bool {
	switch o := i.(type) {
	case ExtractFromHTMLConfiguration:
		return reflect.DeepEqual(o, ExtractFromHTMLConfiguration{})
	case FormatURLConfiguration:
		return reflect.DeepEqual(o, FormatURLConfiguration{})
	case []ScrapeItemConfiguration:
		return len(o) == 0
	case []ScrapeURLConfiguration:
		return len(o) == 0
	case map[string]ExtractFromHTMLConfiguration:
		return len(o) == 0
	case map[string]struct{}:
		return len(o) == 0
	default:
		return false
	}
}
