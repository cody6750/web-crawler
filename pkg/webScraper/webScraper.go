package webcrawler

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
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
	BlackListedURLPaths map[string]struct{}
	RootURL             string
	ScraperNumber       int
	// Queue         chan []string
	Stop        chan struct{}
	WaitGroup   sync.WaitGroup
	HeaderKey   string
	HeaderValue string
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
	webScraper := &WebScraper{}
	return webScraper
}

//Scrape ..
func (w *WebScraper) Scrape(url *URL, scrapeItemConfiguration []ScrapeItemConfiguration, scrapeURLConfiguration ...ScrapeURLConfiguration) (*ScrapeResposne, error) {
	var (
		extractedURL       string
		extractedURLObject *URL
		ExtractedURLs      []*URL
		ExtractedItems     []*Item
		urlTagsToCheck     map[string]bool
		itemTagsToCheck    map[string]bool
		URLsToCheck        map[string]bool
	)
	// log.Print(url)
	// log.Print(scrapeItemConfiguration)
	// log.Print(scrapeURLConfiguration)
	//log.Printf("Scraping link: %v", url)
	URLsToCheck = make(map[string]bool)
	response := ConnectToWebsite(url.CurrentURL, w.HeaderKey, w.HeaderValue).Body
	// body, _ := ioutil.ReadAll(response)
	// log.Print(string(body))
	// writeToFile(string(body))
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
			// log.Print("Start: " + t.String())
			if IsEmpty(scrapeURLConfiguration) {
				//TODO: Replace ExtractedURL with a channel
				extractedURL, err = ExtractURL(t, URLsToCheck)
				log.Print("extracting all url")
				if err != nil {
					continue
				}
			} else {
				extractedURL, err = ExtractURLWithScrapURLConfiguration(t, URLsToCheck, urlTagsToCheck, scrapeURLConfiguration)
				if err != nil {
					continue
				}
			}
			if extractedURL != "" && !IsEmpty(w.BlackListedURLPaths) {
				isBlackListedURLPath := w.isBlackListedURLPath(extractedURL)
				if !isBlackListedURLPath {
					extractedURLObject = &URL{CurrentURL: extractedURL, ParentURL: url.CurrentURL, RootURL: w.RootURL, CurrentDepth: url.CurrentDepth + 1, MaxDepth: url.MaxDepth}
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
			// t := z.Token()
			// log.Print("Error: " + t.String())
			return &ScrapeResposne{RootURL: w.RootURL, ExtractedURLs: ExtractedURLs, ExtractedItem: ExtractedItems}, nil
		}
	}
}

func writeToFile(body string) {
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

func (w *WebScraper) isBlackListedURLPath(url string) bool {
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
		if _, exist := w.BlackListedURLPaths[urlToCheck]; exist {
			log.Printf("URL %s is blacklisted", url)
			return true
		}
		urlToCheck += "/"
		if _, exist := w.BlackListedURLPaths[urlToCheck]; exist {
			log.Printf("URL %s is blacklisted", url)
			return true
		}
	}
	return false
}
