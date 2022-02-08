package webcrawler

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

const (
	hrefAttribute string = "href"
)

//URL ...
type URL struct {
	RootURL      string
	ParentURL    string
	CurrentURL   string
	CurrentDepth int
	MaxDepth     int
}

//ScrapeURLConfig ...
type ScrapeURLConfig struct {
	Name                   string                 `json:"Name"`
	ExtractFromTokenConfig ExtractFromTokenConfig `json:"ExtractFromTokenConfig"`
	FormatURLConfig        FormatURLConfig        `json:"FormatURLConfiguration"`
}

//ExtractURL ...
func ExtractURL(t html.Token, extractedUrls map[string]bool) string {
	url, _ := extractURLFromToken(t)
	if url != "" && !isDuplicateURL(url, extractedUrls) {
		return url
	}
	return ""
}

//ExtractURLWithScrapURLConfig ...
func ExtractURLWithScrapURLConfig(t html.Token, urlsToCheck map[string]bool, tagsToCheck map[string]bool, scrapeURLConfigs []ScrapeURLConfig) (string, error) {
	var url string
	for _, scrapeURLConfig := range scrapeURLConfigs {
		if !IsEmpty(scrapeURLConfig.ExtractFromTokenConfig) {
			if _, exist := tagsToCheck[t.Data]; exist {
				url, _ = extractURLFromTokenUsingConfig(t, scrapeURLConfig.ExtractFromTokenConfig)
			}
		} else {
			url, _ = extractURLFromToken(t)
		}
		if url == "" {
			continue
		}
		if !IsEmpty(scrapeURLConfig.FormatURLConfig) {
			formatedURL := formatURL(url, scrapeURLConfig.FormatURLConfig)
			if formatedURL == "" {
				continue
			}
			if !isDuplicateURL(formatedURL, urlsToCheck) && url != "" {
				return formatedURL, nil

			}
		} else {
			return url, nil
		}
	}
	return url, fmt.Errorf("Unable to extract url with scrap url config")
}
func extractURLFromTokenUsingConfig(token html.Token, urlConfig ExtractFromTokenConfig) (string, error) {
	value, err := extractAttributeValue(token, urlConfig.Attribute)
	if err != nil {
		return value, err
	}
	if strings.Contains(value, urlConfig.AttributeValue) {
		hrefValue, _ := extractAttributeValue(token, hrefAttribute)
		return hrefValue, nil
	}
	return value, fmt.Errorf("unable to extract url from token using the url configuration, %v retrived", value)
}

func extractURLFromToken(token html.Token) (string, error) {
	hrefValue, err := extractAttributeValue(token, hrefAttribute)
	if err != nil {
		return hrefValue, err
	}
	return hrefValue, nil
}

func isDuplicateURL(url string, urlsToCheck map[string]bool) bool {
	if _, exist := urlsToCheck[url]; exist {
		return true
	}
	urlsToCheck[url] = true
	return false
}
