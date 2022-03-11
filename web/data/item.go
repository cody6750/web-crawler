package data

import (
	"encoding/json"
	"net/http"

	webcrawler "github.com/cody6750/web-crawler/pkg"
	webscraper "github.com/cody6750/web-crawler/pkg/webScraper"
	"github.com/sirupsen/logrus"
)

// Item ...
type Item struct {
	Name                      string            `json:"name"`
	Description               string            `json:"description"`
	Price                     float32           `json:"price"`
	URL                       string            `json:"url"`
	ParentURL                 string            `json:"-"`
	AdditionalItemInformation map[string]string `json:"additional_Item_Information"`
	CreatedOn                 string            `json:"-"`
	UpdatedOn                 string            `json:"-"`
	DeletedOn                 string            `json:"-"`
}

//Payload ...
type Payload struct {
	ScrapeItemConfiguration []webscraper.ScrapeItemConfig `json:"ScrapeItemConfiguration"`
	ScrapeURLConfiguration  []webscraper.ScrapeURLConfig  `json:"ScrapeURLConfiguration"`
	RootURL                 string                        `json:"RootURL"`
}

// DecodeToPayload used to decode web crawler response into a usuable struct that the web crawler server
// can parse and execute functions with.
func DecodeToPayload(r *http.Request) (Payload, error) {
	payload := Payload{}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		return Payload{}, err
	}
	return payload, nil
}

// GetItem executes the crawl function within the web crawler. Returns a http response with
// the JSON response as the response body.
func GetItem(crawler *webcrawler.WebCrawler, logger *logrus.Logger, url string, itemsToget []webscraper.ScrapeItemConfig, ScrapeURLConfiguration ...webscraper.ScrapeURLConfig) (*webcrawler.Response, error) {
	crawler.Logger = logger
	response, err := crawler.Crawl(url, itemsToget, ScrapeURLConfiguration...)
	if err != nil {
		logger.WithError(err).Errorf("Failed to get item")
		return response, err
	}
	return response, nil
}
