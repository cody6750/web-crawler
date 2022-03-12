package webcrawler

import (
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
	"golang.org/x/net/html"
)

var (
	err error
)

// WebScraper represents all dependencies required to initialize the web scraper.
type WebScraper struct {

	// Logger used to log.
	Logger *logrus.Logger

	// BlackListedURLPaths used to check if url is blacklisted
	BlackListedURLPaths map[string]struct{}

	// ScraperNumber used to indentify web scraper worker
	ScraperNumber int

	// Stop serves as a signal reciever, used to stop the execution of the web scraper.
	Stop chan struct{}

	// WaitGroup used to wait for channels in the web scraper.
	WaitGroup sync.WaitGroup

	// HeaderKey Used to set http header
	HeaderKey string

	// HeaderKey Used to set http header
	HeaderValue string
}

//Response represents the response the web scraper returns to the web cralwer.
type Response struct {
	RootURL       string
	ExtractedItem []*Item
	ExtractedURLs []*URL
}

//New initializes a web scraper with default options
func New() *WebScraper {
	ws := &WebScraper{}
	ws.Logger = logrus.New()
	return ws
}

//Scrape serves as the main function for the web scraper. Given a url, get the html contents of the url and parse
// the html content for urls and items. It parses the html content by generating tokens for each html element. Tags
// and attributes are extracted for each token,and are used to extract the url and items based on the config parameters.
func (ws *WebScraper) Scrape(u *URL, itemsToGet []ScrapeItemConfig, urlsToGet ...ScrapeURLConfig) (*Response, error) {
	var (
		url             string
		urls            []*URL
		items           []*Item
		itemTagsToCheck map[string]bool
		urlTagsToCheck  map[string]bool
		urlsToCheck     map[string]bool = make(map[string]bool)
	)
	response, err := ConnectToWebsite(u.CurrentURL, ws.HeaderKey, ws.HeaderValue)
	if err != nil {
		return &Response{}, err
	}
	body := response.Body
	if !IsEmpty(urlsToGet) {
		urlTagsToCheck = ws.generateTagsToCheckMap(urlsToGet)
	}
	if !IsEmpty(itemsToGet) {
		itemTagsToCheck = ws.generateTagsToCheckMap(itemsToGet)
	}
	defer body.Close()
	// Parse HTML response by turning it into Tokens
	z := html.NewTokenizer(body)
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
				url = ExtractURL(t, urlsToCheck)
			} else {
				url, _ = ExtractURLWithScrapURLConfig(t, urlsToCheck, urlTagsToCheck, urlsToGet)
			}

			if url != "" {
				if isBlackListedURLPath := ws.isBlackListedURLPath(url); !isBlackListedURLPath {
					urls = append(urls, &URL{CurrentURL: url, ParentURL: u.CurrentURL, RootURL: u.RootURL, CurrentDepth: u.CurrentDepth + 1, MaxDepth: u.MaxDepth})
				}
			}

			if !IsEmpty(itemsToGet) {
				item, err := ExtractItemWithScrapItemConfig(t, z, itemTagsToCheck, itemsToGet)
				if err != nil || IsEmpty(item) {
					continue
				}
				item.URL = u
				items = append(items, &item)
			}

			// This is our break statement
		case tt == html.ErrorToken:
			return &Response{RootURL: u.RootURL, ExtractedURLs: urls, ExtractedItem: items}, nil
		}
	}
}

// generateTagsToCheckMap generates a map of tags to check, given the scrape tag configuration. The map is used for both
// url and item scraping. The map is used to check whether or not the html element should be used to extract from.
func (ws *WebScraper) generateTagsToCheckMap(t interface{}) map[string]bool {
	switch t := t.(type) {
	case []ScrapeURLConfig:
		var urlTagsToCheck = make(map[string]bool)
		for _, ScrapeURLConfig := range t {
			if !IsEmpty(ScrapeURLConfig.ExtractFromTokenConfig) {
				urlTagsToCheck[ScrapeURLConfig.ExtractFromTokenConfig.Tag] = true
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

// isBlackListedURLPath breaks down a url path and checks if it is blacklisted.
func (ws *WebScraper) isBlackListedURLPath(url string) bool {
	var urlToCheck string
	splitURLPath := strings.SplitN(url, "/", 4)
	if len(splitURLPath) < 4 {
		return false
	}
	urlPath := "/" + splitURLPath[3]
	splitURLPath = strings.Split(urlPath, "/")
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
