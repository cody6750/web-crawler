package webcrawler

import (
	"errors"

	"golang.org/x/net/html"
)

var (
	errFormatURL                            error = errors.New("")
	errEmptyParameter                       error = errors.New("")
	errExtractURLFromHTMLUsingConfiguration error = errors.New("")
	errExtractURLFromHTML                   error = errors.New("")
)

//WebScraper ...
type WebScraper struct {
}

//ScrapeConfiguration ...
type ScrapeConfiguration struct {
	ConfigurationName               string
	ExtractURLFromHTMLConfiguration ExtractURLFromHTMLConfiguration
	FormatURLConfiguration          FormatURLConfiguration
}

//ExtractURLFromHTMLConfiguration ...
type ExtractURLFromHTMLConfiguration struct {
	TagToCheck            string
	AttributeToCheck      string
	AttributeValueToCheck string
}

//FormatURLConfiguration ...
type FormatURLConfiguration struct {
	SuffixExist      string
	SuffixToAdd      string
	SuffixToRemove   string
	PrefixToAdd      string
	PrefixExist      string
	PrefixToRemove   string
	ReplaceOldString string
	ReplaceNewString string
}

//New ..
func New() *WebScraper {
	webScraper := &WebScraper{}
	return webScraper
}

//Scrape ..
func (WebScraper) Scrape(url string, scrapeConfiguration ...ScrapeConfiguration) ([]string, error) {
	var (
		ExtractedURLs []string
		TagsToCheck   map[string]bool
		URLsToCheck   map[string]bool
	)

	URLsToCheck = make(map[string]bool)
	response := ConnectToWebsite(url).Body
	defer response.Close()
	// Parse HTML response by turning it into Tokens
	z := html.NewTokenizer(response)

	// If htmlURLConfiguration parameter is provided, create a map filled with tags that will be used to determine if processing is needed. To increase performance
	for _, scrapeConfiguration := range scrapeConfiguration {
		if !isEmptyExtractURLFromHTMLConfiguration(scrapeConfiguration.ExtractURLFromHTMLConfiguration) {
			if len(TagsToCheck) == 0 {
				TagsToCheck = make(map[string]bool)
			}
			TagsToCheck[scrapeConfiguration.ExtractURLFromHTMLConfiguration.TagToCheck] = true
		}
	}
	// This while loop parses through all of the tokens generated for the HTML response.
	for {
		//Iterate through each token
		tt := z.Next()
		// For every token, we check the token type. We parse URL from the start token.
		switch {
		case tt == html.StartTagToken:
			t := z.Token()
			if isEmptyScrapeConfiguration(scrapeConfiguration) {
				//TODO: Replace ExtractedURL with a channel
				_, err := ExtractURL(t, URLsToCheck, ExtractedURLs)
				if err != nil {
					continue
				}
			} else {
				_, err := ExtractURLWithScrapConfiguration(t, URLsToCheck, ExtractedURLs, TagsToCheck, scrapeConfiguration)
				if err != nil {
					continue
				}
			}
			// This is our break statement
		case tt == html.ErrorToken:
			return ExtractedURLs, nil
		}
	}
}

func isEmptyScrapeConfiguration(s []ScrapeConfiguration) bool {
	return len(s) == 0
}
