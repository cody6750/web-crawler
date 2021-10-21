package webscraper

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

const (
	hrefAttribute string = "href"
)

//WebScraper ...
type WebScraper struct {
}

//... ExtractHtmlLinkConfig
type ExtractHtmlLinkConfig struct {
	TagToCheck            string
	AttributeToCheck      string
	AttributeValueToCheck string
}

//New ..
func New() *WebScraper {
	webScraper := &WebScraper{}
	return *&webScraper
}

//Scrape ..
func (WebScraper) Scrape(url string, client *http.Client, htmlLinkConfig ...ExtractHtmlLinkConfig) ([]string, error) {
	var (
		TagsToCheck map[string]bool
		URLsToCheck map[string]bool
	)
	URLsToCheck = make(map[string]bool)
	//GET request to domain for HTML response
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Set header as User-Agent so the server admins don't block our IP address from HTTP requests
	request.Header.Set("User-Agent", "This bot just searches amazon for a product")

	// Make HTTP request
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// Parse HTML response by turning it into Tokens
	z := html.NewTokenizer(response.Body)

	// If htmlLinkConfig parameter is provided, create a map filled with tags that will be used to determine if processing is needed. To increase performance
	if !isEmpty(htmlLinkConfig) {
		TagsToCheck, err = createMapToCheckTags(htmlLinkConfig)
		if err != nil {
			return nil, nil
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
			// If an htmlLinkConfig is provided and the current tags is one of the tags that are required to be checked, determine if the token meets the requirements
			// by checing for the required attribute and attribute value
			if !isEmpty(htmlLinkConfig) {
				if _, tagExist := TagsToCheck[t.Data]; tagExist {
					for _, htmlLinkConfig := range htmlLinkConfig {
						url, _ := extractURLWithHTMLLinkConfig(t, htmlLinkConfig)
						if url != "" && !isDuplicateURL(url, URLsToCheck) {
							log.Printf("URL: %v ", url)
						}
					}
				}
			} else {
				//If an htmlLink is not provided, scrape all href attributes for URLs
				extractURLWithHTMLToken(t)
			}
		// This is our break statement
		case tt == html.ErrorToken:
			return nil, nil
		}
	}
}

func createMapToCheckTags(htmlLinkconfig []ExtractHtmlLinkConfig) (map[string]bool, error) {
	TagsToCheckMap := make(map[string]bool)
	for _, LinkConfig := range htmlLinkconfig {
		if _, exist := TagsToCheckMap[LinkConfig.TagToCheck]; !exist {
			TagsToCheckMap[LinkConfig.TagToCheck] = true
		}
	}
	return TagsToCheckMap, nil
}

func extractURLWithHTMLLinkConfig(token html.Token, htmlLinkConfig ExtractHtmlLinkConfig) (string, error) {
	if htmlLinkConfig.isEmpty() {
		log.Print("is empty")
		return "", errors.New("Empty struct")
	}
	if token.Data == "" {
		log.Print("is empty1")
		return "", errors.New("")
	}
	HTTPAttributeValueFromToken, _ := getHTTPAttributeValueFromToken(token, htmlLinkConfig.AttributeToCheck)
	if strings.Contains(HTTPAttributeValueFromToken, htmlLinkConfig.AttributeValueToCheck) {
		hrefValue, _ := getHTTPAttributeValueFromToken(token, "href")
		return hrefValue, nil
	}
	return "", errors.New("")
}

func extractURLWithHTMLToken(token html.Token) (string, error) {
	attributeValue, error := getHTTPAttributeValueFromToken(token, "href")
	if error != nil {
		return "", error
	}
	if attributeValue == "" {
		return attributeValue, errors.New("TODO")
	}
	return attributeValue, error
}

func getHTTPAttributeValueFromToken(token html.Token, attributeToGet string) (attributeValue string, err error) {
	if attributeToGet == "" {
		return attributeToGet, errors.New("TODO")
	}
	for _, a := range token.Attr {
		if a.Key == attributeToGet {
			attributeValue = a.Val
			return attributeValue, nil
		}
	}
	if attributeValue == "" {
		return attributeValue, errors.New("TODO")
	}
	return attributeValue, nil
}

func (e ExtractHtmlLinkConfig) isEmpty() bool {
	if e != (ExtractHtmlLinkConfig{}) {
		return false
	}
	return true
}

func isEmpty(e []ExtractHtmlLinkConfig) bool {
	if len(e) == 0 {
		return true
	}
	return false
}

func isDuplicateURL(url string, URLsToCheck map[string]bool) bool {
	if _, urlExist := URLsToCheck[url]; urlExist {
		return true
	}
	URLsToCheck[url] = true
	return false
}
