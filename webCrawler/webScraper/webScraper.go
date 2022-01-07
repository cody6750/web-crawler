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

//ScrapeResposne ...
type ScrapeResposne struct {
	RootURL       string
	ExtractedItem []*Item
	ExtractedURLs []*URL
}

//WebScraper ...
type WebScraper struct {
	RootURL       string
	ScraperNumber int
	Queue         chan []string
	Stop          chan struct{}
	WaitGroup     sync.WaitGroup
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
func (w WebScraper) Scrape(url *URL, scrapeItemConfiguration []ScrapeItemConfiguration, scrapeURLConfiguration ...ScrapeURLConfiguration) (*ScrapeResposne, error) {
	var (
		extractedURL       string
		extractedURLObject *URL
		ExtractedURLs      []*URL
		ExtractedItems     []*Item
		urlTagsToCheck     map[string]bool
		itemTagsToCheck    map[string]bool
		URLsToCheck        map[string]bool
	)
	//log.Printf("Scraping link: %v", url)
	URLsToCheck = make(map[string]bool)
	response := ConnectToWebsite(url.CurrentURL).Body
	if !IsEmpty(scrapeURLConfiguration) {
		urlTagsToCheck, err = generateURLTagsToCheckMap(urlTagsToCheck, scrapeURLConfiguration)
		if err != nil {
			return nil, errors.New("Failed to generate url tags")
		}
	}
	if !IsEmpty(scrapeItemConfiguration) {
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
			if IsEmpty(scrapeURLConfiguration) {
				//TODO: Replace ExtractedURL with a channel
				extractedURL, err = ExtractURL(t, URLsToCheck)
				if err != nil {
					continue
				}
			} else {
				extractedURL, err = ExtractURLWithScrapURLConfiguration(t, URLsToCheck, urlTagsToCheck, scrapeURLConfiguration)
				if err != nil {
					continue
				}
			}
			if extractedURL != "" {
				extractedURLObject = &URL{CurrentURL: extractedURL, ParentURL: url.CurrentURL, RootURL: w.RootURL, CurrentDepth: url.CurrentDepth + 1, MaxDepth: url.MaxDepth}
				ExtractedURLs = append(ExtractedURLs, extractedURLObject)
			}
			if !IsEmpty(scrapeItemConfiguration) {
				extractedItem, err := ExtractItemWithScrapItemConfiguration(t, z, itemTagsToCheck, scrapeItemConfiguration)
				if err != nil {
					continue
				}
				extractedItem.URL = extractedURLObject
				ExtractedItems = append(ExtractedItems, &extractedItem)
			}

			// This is our break statement
		case tt == html.ErrorToken:
			return &ScrapeResposne{RootURL: w.RootURL, ExtractedURLs: ExtractedURLs, ExtractedItem: ExtractedItems}, nil
		}
	}
}

func generateItemTagsToCheckMap(itemTagsToCheck map[string]bool, scrapeItemConfiguration []ScrapeItemConfiguration) (map[string]bool, error) {
	if IsEmpty(scrapeItemConfiguration) {
		return nil, errors.New("Item is empty")
	}
	// If item parameter is provided, create a map filled with tags that will be used to determine if processing is needed. To increase performance
	for _, item := range scrapeItemConfiguration {
		if !IsEmpty(item.ItemToGet) {
			if len(itemTagsToCheck) == 0 {
				itemTagsToCheck = make(map[string]bool)
			}
			itemTagsToCheck[item.ItemToGet.Tag] = true
		}
	}
	return itemTagsToCheck, nil
}

func generateURLTagsToCheckMap(urlTagsToCheck map[string]bool, scrapeURLConfiguration []ScrapeURLConfiguration) (map[string]bool, error) {
	if IsEmpty(scrapeURLConfiguration) {
		return nil, errors.New("Scrap configuration is empty")
	}
	// If htmlURLConfiguration parameter is provided, create a map filled with tags that will be used to determine if processing is needed. To increase performance
	for _, scrapeURLConfiguration := range scrapeURLConfiguration {
		if !IsEmpty(scrapeURLConfiguration.ExtractFromHTMLConfiguration) {
			if len(urlTagsToCheck) == 0 {
				urlTagsToCheck = make(map[string]bool)
			}
			urlTagsToCheck[scrapeURLConfiguration.ExtractFromHTMLConfiguration.Tag] = true
		}
	}
	return urlTagsToCheck, nil
}
