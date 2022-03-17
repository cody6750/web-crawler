package webcrawler

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

const (
	hrefAttribute string = "href"
)

//URL represents extracted url
type URL struct {
	RootURL      string
	ParentURL    string
	CurrentURL   string
	CurrentDepth int
	MaxDepth     int
}

//ScrapeURLConfig configuration used to extract url from html token
type ScrapeURLConfig struct {
	Name                   string                 `json:"Name"`
	ExtractFromTokenConfig ExtractFromTokenConfig `json:"ExtractFromTokenConfig"`
	FormatURLConfig        FormatURLConfig        `json:"FormatURLConfiguration"`
}

//ExtractURL extracts url from html token. Checks for duplicates
func ExtractURL(t html.Token, extractedUrls map[string]bool) string {
	url, _ := extractURLFromToken(t)
	if url != "" && !isDuplicateURL(url, extractedUrls) && isURL(url) {
		return url
	}
	return ""
}

// ExtractURLWithScrapURLConfig extracts url from html token using a list of scrape url config which allows
// for selective extraction.
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
		if url == "" || !isURL(url) {
			continue
		}
		if !IsEmpty(scrapeURLConfig.FormatURLConfig) {
			formatedURL := formatURL(url, scrapeURLConfig.FormatURLConfig)
			if formatedURL == "" || !isURL(formatedURL) {
				continue
			}
			if !isDuplicateURL(formatedURL, urlsToCheck) && url != "" {
				return formatedURL, nil

			}
		} else {
			if !isURL(url) {
				continue
			}
			return url, nil
		}
	}
	return url, fmt.Errorf("Unable to extract url with scrap url config")
}

// extractURLFromTokenUsingConfig extracts url from html token using a scrape url config which allows
// for selective extraction.
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

// extractURLFromTokenUsingConfig extracts url from html token which allows for selective extraction.
func extractURLFromToken(token html.Token) (string, error) {
	hrefValue, err := extractAttributeValue(token, hrefAttribute)
	if err != nil {
		return hrefValue, err
	}
	return hrefValue, nil
}

// isDuplicateURL checks for duplicate urls.
func isDuplicateURL(url string, urlsToCheck map[string]bool) bool {
	if _, exist := urlsToCheck[url]; exist {
		return true
	}
	urlsToCheck[url] = true
	return false
}
