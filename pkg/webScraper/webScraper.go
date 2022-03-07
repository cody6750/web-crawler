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
					urls = append(urls, &URL{CurrentURL: url, ParentURL: u.CurrentURL, RootURL: ws.RootURL, CurrentDepth: u.CurrentDepth + 1, MaxDepth: u.MaxDepth})
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
			return &Response{RootURL: ws.RootURL, ExtractedURLs: urls, ExtractedItem: items}, nil
		}
	}
}

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
