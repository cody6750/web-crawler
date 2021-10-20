package webscraper

import (
	"log"
	"net/http"
)

//WebScraper ...
type WebScraper struct {
}

//New ..
func (WebScraper) New() *WebScraper {
	webScraper := &WebScraper{}
	return *&webScraper
}

//Scrape ..
func (WebScraper) Scrape(url string, client *http.Client) {
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
}
