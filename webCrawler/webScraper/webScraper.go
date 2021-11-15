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

var (
	errFormatURL                            error = errors.New("")
	errEmptyParameter                       error = errors.New("")
	errExtractURLFromHTMLUsingConfiguration error = errors.New("")
	errExtractURLFromHTML                   error = errors.New("")
)

//WebScraper ...
type WebScraper struct {
}

//ScrapeConfiguration ...
type ScrapeConfiguration struct {
	ConfigurationName               string
	ExtractURLFromHTMLConfiguration ExtractURLFromHTMLConfiguration
	FormatURLConfiguration          FormatURLConfiguration
}

//ExtractURLFromHTMLConfiguration ...
type ExtractURLFromHTMLConfiguration struct {
	TagToCheck            string
	AttributeToCheck      string
	AttributeValueToCheck string
}

//FormatURLConfiguration ...
type FormatURLConfiguration struct {
	SuffixToAdd      string
	SuffixToRemove   string
	PrefixToAdd      string
	PrefixToRemove   string
	ReplaceOldString string
	ReplaceNewString string
}

//New ..
func New() *WebScraper {
	webScraper := &WebScraper{}
	return *&webScraper
}

//Scrape ..
func (WebScraper) Scrape(url string, client *http.Client, scrapeConfiguration ...ScrapeConfiguration) ([]string, error) {
	var (
		ExtractedURLs []string
		TagsToCheck   map[string]bool
		URLsToCheck   map[string]bool
		ExtractedURL  string
		counter       int
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

	// If htmlURLConfiguration parameter is provided, create a map filled with tags that will be used to determine if processing is needed. To increase performance
	for _, scrapeConfiguration := range scrapeConfiguration {
		if !isEmptyFormatURLConfiguration(scrapeConfiguration.FormatURLConfiguration) {
			if len(TagsToCheck) == 0 {
				TagsToCheck = make(map[string]bool)
			}
			TagsToCheck[scrapeConfiguration.ExtractURLFromHTMLConfiguration.TagToCheck] = true
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
			if isEmptyScrapeConfiguration(scrapeConfiguration) {
				ExtractedURL, _ = extractURLFromHTML(t)
				if ExtractedURL != "" && !isDuplicateURL(ExtractedURL, URLsToCheck) {
					log.Default().Printf("Extracted url: %v", ExtractedURL)
					ExtractedURLs = append(ExtractedURLs, ExtractedURL)
				}
			} else {
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

						}
						if !isDuplicateURL(formatedURL, URLsToCheck) {
							ExtractedURLs = append(ExtractedURLs, formatedURL)
							log.Printf("%v. Formated url: %v", counter, formatedURL)
							counter++
						}
					} else {
						log.Default().Printf("Extracted url: %v", ExtractedURL)
						ExtractedURLs = append(ExtractedURLs, ExtractedURL)
					}
					ExtractedURL = ""
				}
			}
			// This is our break statement
		case tt == html.ErrorToken:
			return ExtractedURLs, nil
		}
	}
}

func extractURLFromHTMLUsingConfiguration(token html.Token, urlConfig ExtractURLFromHTMLConfiguration) (string, error) {
	if isEmptyExtractURLFromHTMLConfiguration(urlConfig) {
		log.Print("is empty")
		return "", errExtractURLFromHTMLUsingConfiguration
	}
	HTTPAttributeValueFromToken, _ := getHTTPAttributeValueFromToken(token, urlConfig.AttributeToCheck)
	if strings.Contains(HTTPAttributeValueFromToken, urlConfig.AttributeValueToCheck) {
		//log.Printf("Got attribute %v, Value %v", urlConfig.AttributeToCheck, HTTPAttributeValueFromToken)
		hrefValue, _ := getHTTPAttributeValueFromToken(token, "href")
		return hrefValue, nil
	}
	return "", errExtractURLFromHTMLUsingConfiguration
}

func extractURLFromHTML(token html.Token) (string, error) {
	if isEmptyToken(token) {
		log.Print("is empty1")
		return "", errExtractURLFromHTML
	}
	hrefValue, error := getHTTPAttributeValueFromToken(token, "href")
	if error != nil {
		return "", errExtractURLFromHTML
	}
	if hrefValue == "" {
		return hrefValue, errExtractURLFromHTML
	}
	return hrefValue, nil
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

func formatURL(url string, formatURLConfig FormatURLConfiguration) (string, error) {
	if strings.Contains(url, formatURLConfig.PrefixToRemove) && strings.Contains(url, formatURLConfig.SuffixToRemove) && strings.Contains(url, formatURLConfig.ReplaceOldString) {
		if formatURLConfig.ReplaceOldString != "" && formatURLConfig.ReplaceNewString != "" {
			url = strings.ReplaceAll(url, formatURLConfig.ReplaceOldString, formatURLConfig.ReplaceNewString)
		}
		if strings.HasPrefix(url, formatURLConfig.PrefixToRemove) {
			url = strings.TrimPrefix(url, formatURLConfig.PrefixToRemove)
		}
		if strings.HasSuffix(url, formatURLConfig.SuffixToRemove) {
			url = strings.TrimSuffix(url, formatURLConfig.SuffixToRemove)
		}
		if formatURLConfig.PrefixToAdd != "" {
			url = formatURLConfig.PrefixToAdd + url
		}
		if formatURLConfig.SuffixToAdd != "" {
			url = url + formatURLConfig.SuffixToAdd
		}
	} else {
		return "", errFormatURL
	}

	return url, nil
}

func isEmptyToken(token html.Token) bool {
	if token.Data == "" && len(token.Attr) == 0 {
		return true
	}
	return false
}
func isEmptyFormatURLConfiguration(formatURLConfiguration FormatURLConfiguration) bool {
	if formatURLConfiguration == (FormatURLConfiguration{}) {
		return true
	}
	return false
}

func isEmptyExtractURLFromHTMLConfiguration(extractURLFromHTMLConfiguration ExtractURLFromHTMLConfiguration) bool {
	if extractURLFromHTMLConfiguration == (ExtractURLFromHTMLConfiguration{}) {
		return true
	}
	return false
}

func isEmptyScrapeConfiguration(s []ScrapeConfiguration) bool {
	if len(s) == 0 {
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
