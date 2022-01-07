package webcrawler

//IsEmpty ...
func IsEmpty(i interface{}) bool {
	switch o := i.(type) {
	case ExtractFromHTMLConfiguration:
		return o == (ExtractFromHTMLConfiguration{})
	case FormatURLConfiguration:
		return o == (FormatURLConfiguration{})
	case []ScrapeItemConfiguration:
		return len(o) == 0
	case []ScrapeURLConfiguration:
		return len(o) == 0
	case map[string]ExtractFromHTMLConfiguration:
		return len(o) == 0
	default:
		return false
	}
}
