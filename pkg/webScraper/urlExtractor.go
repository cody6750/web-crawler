package webcrawler

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

const (
	hrefAttribute string = "href"
)

//ScrapeURLConfig ...
type ScrapeURLConfig struct {
	Name                         string                       `json:"Name"`
	ExtractFromHTMLConfiguration ExtractFromHTMLConfiguration `json:"ExtractFromHTMLConfiguration"`
	FormatURLConfiguration       FormatURLConfiguration       `json:"FormatURLConfiguration"`
}

//ExtractURL ...
func ExtractURL(t html.Token, URLsToCheck map[string]bool) (string, error) {
	ExtractedURL, _ := extractURLFromHTML(t)
	if ExtractedURL != "" && !isDuplicateURL(ExtractedURL, URLsToCheck) {
		return ExtractedURL, nil
	}
	return ExtractedURL, nil
}

//ExtractURLWithScrapURLConfig ...
func ExtractURLWithScrapURLConfig(t html.Token, URLsToCheck map[string]bool, TagsToCheck map[string]bool, scrapeURLConfiguration []ScrapeURLConfig) (string, error) {
	var ExtractedURL string
	for _, scrapeURLConfiguration := range scrapeURLConfiguration {
		if !IsEmpty(scrapeURLConfiguration.ExtractFromHTMLConfiguration) {
			if _, tagExist := TagsToCheck[t.Data]; tagExist {
				ExtractedURL, _ = extractURLFromHTMLUsingConfiguration(t, scrapeURLConfiguration.ExtractFromHTMLConfiguration)
			}
		} else {
			ExtractedURL, _ = extractURLFromHTML(t)
		}
		if ExtractedURL == "" {
			continue
		}
		if !IsEmpty(scrapeURLConfiguration.FormatURLConfiguration) {
			formatedURL := formatURL(ExtractedURL, scrapeURLConfiguration.FormatURLConfiguration)
			if formatedURL == "" {
				continue
			}
			if !isDuplicateURL(formatedURL, URLsToCheck) && ExtractedURL != "" {
				return formatedURL, nil

			}
		} else {
			return ExtractedURL, nil
		}
	}
	return "", nil
}
func extractURLFromHTMLUsingConfiguration(token html.Token, urlConfig ExtractFromHTMLConfiguration) (string, error) {
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

func extractURLFromHTML(token html.Token) (string, error) {
	hrefValue, error := extractAttributeValue(token, hrefAttribute)
	if error != nil {
		return "", errExtractURLFromHTML
	}
	if hrefValue == "" {
		return hrefValue, errExtractURLFromHTML
	}
	return hrefValue, nil
}

func isDuplicateURL(url string, URLsToCheck map[string]bool) bool {
	if _, urlExist := URLsToCheck[url]; urlExist {
		return true
	}
	URLsToCheck[url] = true
	return false
}
