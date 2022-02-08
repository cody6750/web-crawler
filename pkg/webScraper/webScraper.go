package webcrawler

import (
	"errors"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
	"golang.org/x/net/html"
)

var (
	errExtractURLFromHTML error = errors.New("")
	err                   error
)

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

//Response ...
type Response struct {
	RootURL       string
	ExtractedItem []*Item
	ExtractedURLs []*URL
}

//New ..
func New() *WebScraper {
	ws := &WebScraper{}
	ws.Logger = logrus.New()
	return ws
}

//Scrape ..
func (ws *WebScraper) Scrape(url *URL, itemsToGet []ScrapeItemConfig, urlsToGet ...ScrapeURLConfig) (*Response, error) {
	var (
		extractedURL       string
		extractedURLObject *URL
		ExtractedURLs      []*URL
		ExtractedItems     []*Item
		itemTagsToCheck    map[string]bool
		urlTagsToCheck     map[string]bool
		URLsToCheck        map[string]bool = make(map[string]bool)
	)
	response := ConnectToWebsite(url.CurrentURL, ws.HeaderKey, ws.HeaderValue).Body
	if !IsEmpty(urlsToGet) {
		urlTagsToCheck = ws.generateTagsToCheckMap(urlsToGet)
	}
	if !IsEmpty(itemsToGet) {
		itemTagsToCheck = ws.generateTagsToCheckMap(itemsToGet)
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
			if IsEmpty(urlsToGet) {
				//TODO: Replace ExtractedURL with a channel
				extractedURL, err = ExtractURL(t, URLsToCheck)
				if err != nil {
					continue
				}
			} else {
				extractedURL, err = ExtractURLWithScrapURLConfig(t, URLsToCheck, urlTagsToCheck, urlsToGet)
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

			if !IsEmpty(itemsToGet) {
				extractedItem, err := ExtractItemWithScrapItemConfig(t, z, itemTagsToCheck, itemsToGet)
				if err != nil {
					continue
				}
				extractedItem.URL = url
				//extractedItem.printJSON()
				ExtractedItems = append(ExtractedItems, &extractedItem)
			}

			// This is our break statement
		case tt == html.ErrorToken:
			return &Response{RootURL: ws.RootURL, ExtractedURLs: ExtractedURLs, ExtractedItem: ExtractedItems}, nil
		}
	}
}

// func generateItemTagsToCheckMap(scrapeItemConfiguration []ScrapeItemConfiguration) map[string]bool {
// 	// If item parameter is provided, create a map filled with tags that will be used to determine if processing is needed. To increase performance
// 	var itemTagsToCheck = make(map[string]bool)
// 	for _, item := range scrapeItemConfiguration {
// 		if !IsEmpty(item.ItemToGet) {
// 			itemTagsToCheck[item.ItemToGet.Tag] = true
// 		}
// 	}
// 	return itemTagsToCheck
// }
func (ws *WebScraper) generateTagsToCheckMap(t interface{}) map[string]bool {
	switch t := t.(type) {
	case []ScrapeURLConfig:
		var urlTagsToCheck = make(map[string]bool)
		for _, ScrapeURLConfig := range t {
			if !IsEmpty(ScrapeURLConfig.ExtractFromHTMLConfiguration) {
				urlTagsToCheck[ScrapeURLConfig.ExtractFromHTMLConfiguration.Tag] = true
			}
		}
		return urlTagsToCheck
	case []ScrapeItemConfig:
		var itemTagsToCheck = make(map[string]bool)
		for _, item := range t {
			if !IsEmpty(item.ItemToGet) {
				itemTagsToCheck[item.ItemToGet.Tag] = true
			}
		}
		return itemTagsToCheck
	default:
		ws.Logger.WithField("Type", t).Warn("Unable to generate tags to check map")
	}
	return map[string]bool{}
}

// func generateURLTagsToCheckMap(ScrapeURLConfig []ScrapeURLConfig) map[string]bool {
// 	// If htmlURLConfiguration parameter is provided, create a map filled with tags that will be used to determine if processing is needed. To increase performance
// 	var urlTagsToCheck = make(map[string]bool)
// 	for _, ScrapeURLConfig := range ScrapeURLConfig {
// 		if !IsEmpty(ScrapeURLConfig.ExtractFromHTMLConfiguration) {
// 			urlTagsToCheck[ScrapeURLConfig.ExtractFromHTMLConfiguration.Tag] = true
// 		}
// 	}
// 	return urlTagsToCheck
// }

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
