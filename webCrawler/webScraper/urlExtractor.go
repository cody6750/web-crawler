package webcrawler

import (
	"errors"
	"log"
	"strings"

	"golang.org/x/net/html"
)

const (
	hrefAttribute string = "href"
)

//ExtractURL ...
func ExtractURL(t html.Token, URLsToCheck map[string]bool, ExtractedURLs []string) (string, error) {
	ExtractedURL, _ := extractURLFromHTML(t)
	if ExtractedURL != "" && !isDuplicateURL(ExtractedURL, URLsToCheck) {
		log.Default().Printf("Extracted url: %v", ExtractedURL)
		ExtractedURLs = append(ExtractedURLs, ExtractedURL)
		return ExtractedURL, nil
	}
	return ExtractedURL, errors.New("")
}

//ExtractURLWithScrapConfiguration ...
func ExtractURLWithScrapConfiguration(t html.Token, URLsToCheck map[string]bool, ExtractedURLs []string, TagsToCheck map[string]bool, scrapeConfiguration []ScrapeConfiguration) (string, error) {
	var ExtractedURL string
	for _, scrapeConfiguration := range scrapeConfiguration {
		if !isEmptyExtractURLFromHTMLConfiguration(scrapeConfiguration.ExtractURLFromHTMLConfiguration) {
			if _, tagExist := TagsToCheck[t.Data]; tagExist {
				ExtractedURL, _ = extractURLFromHTMLUsingConfiguration(t, scrapeConfiguration.ExtractURLFromHTMLConfiguration)
			}
		} else {
			ExtractedURL, _ = extractURLFromHTML(t)
		}
		if ExtractedURL == "" {
			continue
		}
		if !isEmptyFormatURLConfiguration(scrapeConfiguration.FormatURLConfiguration) {
			formatedURL, err := formatURL(ExtractedURL, scrapeConfiguration.FormatURLConfiguration)
			if err != nil {
				continue
			}

			if !isDuplicateURL(formatedURL, URLsToCheck) {
				log.Default().Printf("Formated url: %v", formatedURL)
				ExtractedURLs = append(ExtractedURLs, formatedURL)
				return formatedURL, nil

			}
		} else {
			log.Default().Printf("Extracted url: %v", ExtractedURL)
			ExtractedURLs = append(ExtractedURLs, ExtractedURL)
			return ExtractedURL, nil
		}
	}
	return "", errors.New("")
}
func extractURLFromHTMLUsingConfiguration(token html.Token, urlConfig ExtractURLFromHTMLConfiguration) (string, error) {
	if isEmptyExtractURLFromHTMLConfiguration(urlConfig) {
		log.Print("is empty")
		return "", errExtractURLFromHTMLUsingConfiguration
	}
	HTTPAttributeValueFromToken, _ := getHTTPAttributeValueFromToken(token, urlConfig.AttributeToCheck)
	if strings.Contains(HTTPAttributeValueFromToken, urlConfig.AttributeValueToCheck) {
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

func isEmptyExtractURLFromHTMLConfiguration(extractURLFromHTMLConfiguration ExtractURLFromHTMLConfiguration) bool {
	return extractURLFromHTMLConfiguration == (ExtractURLFromHTMLConfiguration{})
}

func getHTTPAttributeValueFromToken(token html.Token, attributeToGet string) (attributeValue string, err error) {
	if attributeToGet == "" {
		return "empty string", errEmptyParameter
	}
	for _, a := range token.Attr {
		if a.Key == attributeToGet {
			attributeValue = a.Val
			return attributeValue, nil
		}
	}
	if attributeValue == "" {
		return "empty string", errEmptyParameter
	}
	return "does not exist", nil
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
