package webcrawler

import (
	"log"
	"strings"

	"golang.org/x/net/html"
)

const (
	hrefAttribute string = "href"
)

//ExtractURL ...
func ExtractURL(t html.Token, URLsToCheck map[string]bool) (string, error) {

	ExtractedURL, _ := extractURLFromHTML(t)
	if ExtractedURL != "" && !isDuplicateURL(ExtractedURL, URLsToCheck) {
		log.Default().Printf("Extracted url: %v", ExtractedURL)
		return ExtractedURL, nil
	}
	return ExtractedURL, nil
}

//ExtractURLWithScrapURLConfiguration ...
func ExtractURLWithScrapURLConfiguration(t html.Token, URLsToCheck map[string]bool, TagsToCheck map[string]bool, scrapeURLConfiguration []ScrapeURLConfiguration) (string, error) {
	var ExtractedURL string
	for _, scrapeURLConfiguration := range scrapeURLConfiguration {
		if !isEmptyExtractFromHTMLConfiguration(scrapeURLConfiguration.ExtractFromHTMLConfiguration) {
			if _, tagExist := TagsToCheck[t.Data]; tagExist {
				ExtractedURL, _ = extractURLFromHTMLUsingConfiguration(t, scrapeURLConfiguration.ExtractFromHTMLConfiguration)
			}
		} else {
			ExtractedURL, _ = extractURLFromHTML(t)
		}
		if ExtractedURL == "" {
			continue
		}
		if !isEmptyFormatURLConfiguration(scrapeURLConfiguration.FormatURLConfiguration) {
			formatedURL, err := formatURL(ExtractedURL, scrapeURLConfiguration.FormatURLConfiguration)
			if err != nil {
				continue
			}

			if !isDuplicateURL(formatedURL, URLsToCheck) && ExtractedURL != "" {
				//log.Default().Printf("Formated url: %v", formatedURL)
				return formatedURL, nil

			}
		} else {
			log.Default().Printf("Extracted url: %v", ExtractedURL)
			return ExtractedURL, nil
		}
	}
	return "", nil
}
func extractURLFromHTMLUsingConfiguration(token html.Token, urlConfig ExtractFromHTMLConfiguration) (string, error) {
	if isEmptyExtractFromHTMLConfiguration(urlConfig) {
		log.Print("is empty")
		return "", errExtractURLFromHTMLUsingConfiguration
	}
	HTTPAttributeValueFromToken, _ := getHTTPAttributeValueFromToken(token, urlConfig.Attribute)
	if strings.Contains(HTTPAttributeValueFromToken, urlConfig.AttributeValue) {
		//log.Printf("Got attribute %v, Value %v", urlConfig.AttributeToCheck, HTTPAttributeValueFromToken)
		hrefValue, _ := getHTTPAttributeValueFromToken(token, hrefAttribute)
		return hrefValue, nil
	}
	return "", errExtractURLFromHTMLUsingConfiguration
}

func extractURLFromHTML(token html.Token) (string, error) {
	if isEmptyToken(token) {
		log.Print("is empty1")
		return "", errExtractURLFromHTML
	}
	hrefValue, error := getHTTPAttributeValueFromToken(token, hrefAttribute)
	if error != nil {
		return "", errExtractURLFromHTML
	}
	if hrefValue == "" {
		return hrefValue, errExtractURLFromHTML
	}
	return hrefValue, nil
}

func isEmptyExtractFromHTMLConfiguration(extractFromHTMLConfiguration ExtractFromHTMLConfiguration) bool {
	return extractFromHTMLConfiguration == (ExtractFromHTMLConfiguration{})
}

func isEmptyToken(token html.Token) bool {
	if token.Data == "" && len(token.Attr) == 0 {
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
