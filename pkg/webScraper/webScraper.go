package webcrawler

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
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
	Logger              *logrus.Logger
	BlackListedURLPaths map[string]struct{}
	RootURL             string
	ScraperNumber       int
	Stop                chan struct{}
	WaitGroup           sync.WaitGroup
	HeaderKey           string
	HeaderValue         string
}

//ScrapeURLConfiguration ...
type ScrapeURLConfiguration struct {
	Name                         string                       `json:"Name"`
	ExtractFromHTMLConfiguration ExtractFromHTMLConfiguration `json:"ExtractFromHTMLConfiguration"`
	FormatURLConfiguration       FormatURLConfiguration       `json:"FormatURLConfiguration"`
}

//FormatURLConfiguration ...
type FormatURLConfiguration struct {
	SuffixExist      string `json:"SuffixExist"`
	SuffixToAdd      string `json:"SuffixToAdd"`
	SuffixToRemove   string `json:"SuffixToRemove"`
	PrefixToAdd      string `json:"PrefixToAdd"`
	PrefixExist      string `json:"PrefixExist"`
	PrefixToRemove   string `json:"PrefixToRemove"`
	ReplaceOldString string `json:"ReplaceOldString"`
	ReplaceNewString string `json:"ReplaceNewString"`
}

//New ..
func New() *WebScraper {
	ws := &WebScraper{}
	ws.Logger = logrus.New()
	return ws
}

//Scrape ..
func (ws *WebScraper) Scrape(url *URL, scrapeItemConfiguration []ScrapeItemConfiguration, scrapeURLConfiguration ...ScrapeURLConfiguration) (*ScrapeResposne, error) {
	var (
		extractedURL       string
		extractedURLObject *URL
		ExtractedURLs      []*URL
		ExtractedItems     []*Item
		urlTagsToCheck     map[string]bool
		itemTagsToCheck    map[string]bool
		URLsToCheck        map[string]bool
	)
	URLsToCheck = make(map[string]bool)
	response := ConnectToWebsite(url.CurrentURL, ws.HeaderKey, ws.HeaderValue).Body
	if !IsEmpty(scrapeURLConfiguration) {
		urlTagsToCheck, err = generateURLTagsToCheckMap(urlTagsToCheck, scrapeURLConfiguration)
		if err != nil {
			return nil, errors.New("failed to generate url tags")
		}
	}
	if !IsEmpty(scrapeItemConfiguration) {
		itemTagsToCheck, err = generateItemTagsToCheckMap(itemTagsToCheck, scrapeItemConfiguration)
		if err != nil {
			return nil, errors.New("failed to generate item tags")
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
			if extractedURL != "" && !IsEmpty(ws.BlackListedURLPaths) {
				isBlackListedURLPath := ws.isBlackListedURLPath(extractedURL)
				if !isBlackListedURLPath {
					extractedURLObject = &URL{CurrentURL: extractedURL, ParentURL: url.CurrentURL, RootURL: ws.RootURL, CurrentDepth: url.CurrentDepth + 1, MaxDepth: url.MaxDepth}
					ExtractedURLs = append(ExtractedURLs, extractedURLObject)
				}
			}

			if !IsEmpty(scrapeItemConfiguration) {
				extractedItem, err := ExtractItemWithScrapItemConfiguration(t, z, itemTagsToCheck, scrapeItemConfiguration)
				if err != nil {
					continue
				}
				extractedItem.URL = url
				// extractedItem.printJSON()
				ExtractedItems = append(ExtractedItems, &extractedItem)
			}

			// This is our break statement
		case tt == html.ErrorToken:
			return &ScrapeResposne{RootURL: ws.RootURL, ExtractedURLs: ExtractedURLs, ExtractedItem: ExtractedItems}, nil
		}
	}
}

func WriteToFile(body string) {
	f, err := os.Create("data.html")

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	_, err2 := f.WriteString(body)

	if err2 != nil {
		log.Fatal(err2)
	}

	fmt.Println("done")
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
		return nil, errors.New("scrap configuration is empty")
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

func (ws *WebScraper) isBlackListedURLPath(url string) bool {
	var urlToCheck string
	splitURLPAth := strings.SplitN(url, "/", 4)
	if len(splitURLPAth) < 4 {
		return false
	}
	urlPath := "/" + splitURLPAth[3]
	splitURLPath := strings.Split(urlPath, "/")
	for _, splitURL := range splitURLPath {
		urlToCheck += splitURL
		if urlToCheck == "" {
			urlToCheck += "/"
			continue
		}
		if _, exist := ws.BlackListedURLPaths[urlToCheck]; exist {
			return true
		}
		urlToCheck += "/"
		if _, exist := ws.BlackListedURLPaths[urlToCheck]; exist {
			return true
		}
	}
	return false
}
