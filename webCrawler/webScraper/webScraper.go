package webcrawler

import (
	"errors"
	"sync"

	"golang.org/x/net/html"
)

var (
	errFormatURL                            error = errors.New("")
	errEmptyParameter                       error = errors.New("")
	errExtractURLFromHTMLUsingConfiguration error = errors.New("")
	errExtractURLFromHTML                   error = errors.New("")
	err                                     error
)

//WebScraper ...
type WebScraper struct {
	Host          string
	ScraperNumber int
	Queue         chan []string
	Stop          chan struct{}
	Wg            sync.WaitGroup
}

//ScrapeURLConfiguration ...
type ScrapeURLConfiguration struct {
	ConfigurationName            string
	ExtractFromHTMLConfiguration ExtractFromHTMLConfiguration
	FormatURLConfiguration       FormatURLConfiguration
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
func (WebScraper) Scrape(url string, scrapeItemConfiguration []ScrapeItemConfiguration, scrapeURLConfiguration ...ScrapeURLConfiguration) ([]string, error) {
	var (
		ExtractedURLs []string
		//ExtractedItems  []Item
		urlTagsToCheck  map[string]bool
		itemTagsToCheck map[string]bool
		URLsToCheck     map[string]bool
	)
	//log.Printf("Scraping link: %v", url)
	URLsToCheck = make(map[string]bool)
	response := ConnectToWebsite(url).Body
	if !isEmptyScrapeURLConfiguration(scrapeURLConfiguration) {
		urlTagsToCheck, err = generateURLTagsToCheckMap(urlTagsToCheck, scrapeURLConfiguration)
		if err != nil {
			return nil, errors.New("Failed to generate url tags")
		}
	}
	if !isEmptyItem(scrapeItemConfiguration) {
		itemTagsToCheck, err = generateItemTagsToCheckMap(itemTagsToCheck, scrapeItemConfiguration)
		if err != nil {
			return nil, errors.New("Failed to generate item tags")
		}
	}
	defer response.Close()
	// Parse HTML response by turning it into Tokens
	z := html.NewTokenizer(response)
	// This while loop parses through all of the tokens generated for the HTML response.
	for {
		//Iterate through each token
		tt := z.Next()
		// For every token, we check the token type. We parse URL from the start token.
		switch {
		case tt == html.StartTagToken:
			t := z.Token()
			if isEmptyScrapeURLConfiguration(scrapeURLConfiguration) {
				//TODO: Replace ExtractedURL with a channel
				extractedURL, err := ExtractURL(t, URLsToCheck)
				if err != nil {
					continue
				}
				if extractedURL != "" {
					ExtractedURLs = append(ExtractedURLs, extractedURL)
				}
			} else {
				extractedURL, err := ExtractURLWithScrapURLConfiguration(t, URLsToCheck, urlTagsToCheck, scrapeURLConfiguration)
				if err != nil {
					continue
				}
				if extractedURL != "" {
					ExtractedURLs = append(ExtractedURLs, extractedURL)
				}
			}

			if !isEmptyItem(scrapeItemConfiguration) {
				err := ExtractItemWithScrapItemConfiguration(t, z, url, itemTagsToCheck, scrapeItemConfiguration)
				if err != nil {
					continue
				}
				//ExtractedItems = append(ExtractedItems, extractedItem)
			}

			// This is our break statement
		case tt == html.ErrorToken:
			return ExtractedURLs, nil
		}
	}
}

func generateItemTagsToCheckMap(itemTagsToCheck map[string]bool, scrapeItemConfiguration []ScrapeItemConfiguration) (map[string]bool, error) {
	if isEmptyItem(scrapeItemConfiguration) {
		return nil, errors.New("Item is empty")
	}
	// If item parameter is provided, create a map filled with tags that will be used to determine if processing is needed. To increase performance
	for _, item := range scrapeItemConfiguration {
		if !isEmptyExtractFromHTMLConfiguration(item.ItemToGet) {
			if len(itemTagsToCheck) == 0 {
				itemTagsToCheck = make(map[string]bool)
			}
			itemTagsToCheck[item.ItemToGet.Tag] = true
		}
	}
	return itemTagsToCheck, nil
}

func generateURLTagsToCheckMap(urlTagsToCheck map[string]bool, scrapeURLConfiguration []ScrapeURLConfiguration) (map[string]bool, error) {
	if isEmptyScrapeURLConfiguration(scrapeURLConfiguration) {
		return nil, errors.New("Scrap configuration is empty")
	}
	// If htmlURLConfiguration parameter is provided, create a map filled with tags that will be used to determine if processing is needed. To increase performance
	for _, scrapeURLConfiguration := range scrapeURLConfiguration {
		if !isEmptyExtractFromHTMLConfiguration(scrapeURLConfiguration.ExtractFromHTMLConfiguration) {
			if len(urlTagsToCheck) == 0 {
				urlTagsToCheck = make(map[string]bool)
			}
			urlTagsToCheck[scrapeURLConfiguration.ExtractFromHTMLConfiguration.Tag] = true
		}
	}
	return urlTagsToCheck, nil
}

func isEmptyScrapeURLConfiguration(s []ScrapeURLConfiguration) bool {
	return len(s) == 0
}

func isEmptyItem(i []ScrapeItemConfiguration) bool {
	return len(i) == 0
}
