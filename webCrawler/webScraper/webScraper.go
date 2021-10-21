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
	)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("User-Agent", "This bot just searches amazon for a product")

	// Make request
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	z := html.NewTokenizer(response.Body)

	// If htmlLinkConfig is provided, create a map to search for supported Tags with constant time O(1).
	if !isEmpty(htmlLinkConfig) {
		TagsToCheck, err = createMapToCheckTags(htmlLinkConfig)
		if err != nil {
			return nil, nil
		}
	}
	for {
		tt := z.Next()
		switch {
		case tt == html.StartTagToken:
			t := z.Token()
			if _, tagExist := TagsToCheck[t.Data]; !isEmpty(htmlLinkConfig) && tagExist {
				for _, htmlLinkConfig := range htmlLinkConfig {
					extractURLWithHTMLLinkConfig(t, htmlLinkConfig)
				}
			} else {
				extractURLWithHTMLToken(t)
			}
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
		return "", errors.New("Empty struct")
	}
	if token.Data == "" {
		return "", errors.New("")
	}
	if HTTPAttributeValueFromToken, _ := getHTTPAttributeValueFromToken(token, htmlLinkConfig.AttributeToCheck); strings.Contains(HTTPAttributeValueFromToken, htmlLinkConfig.AttributeValueToCheck) {
		hrefValue, _ := getHTTPAttributeValueFromToken(token, "href")
		//log.Printf("HREF URL: %v\n TOKEN ATTR: %v", hrefValue, token.Attr)
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
	log.Printf("HREF URL: %v", attributeValue)
	return attributeValue, error
}

func getHTTPAttributeValueFromToken(token html.Token, attributeToGet string) (attributeValue string, err error) {
	if attributeToGet == "" {
		return attributeToGet, errors.New("TODO")
	}
	for _, a := range token.Attr {
		if a.Key == attributeToGet {
			attributeValue = a.Val
			break
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
