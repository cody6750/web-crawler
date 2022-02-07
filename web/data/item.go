package data

import (
	"encoding/json"
	"io"
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
	ScrapeItemConfiguration []webscraper.ScrapeItemConfiguration `json:"ScrapeItemConfiguration"`
	ScrapeURLConfiguration  []webscraper.ScrapeURLConfiguration  `json:"ScrapeURLConfiguration"`
	RootURL                 string                               `json:"RootURL"`
}

// Items ...
type Items []*Item

// DecodeToPayload ...
func DecodeToPayload(r *http.Request) (*Payload, error) {
	payload := &Payload{}
	err := json.NewDecoder(r.Body).Decode(payload)
	if err != nil {
		return nil, err
	}
	return payload, nil
}

// GetItem ...
func GetItem(logger *logrus.Logger, url string, itemsToget []webscraper.ScrapeItemConfiguration, ScrapeURLConfiguration ...webscraper.ScrapeURLConfiguration) ([]*webscraper.ScrapeResposne, error) {
	crawler := webcrawler.NewCrawler()
	crawler.Logger = logger
	response, err := crawler.Crawl(url, itemsToget, ScrapeURLConfiguration...)
	if err != nil {
		logger.WithError(err).Errorf("Failed to get item")
		return response, err
	}
	return response, nil
}

// ToJSON ...
func ToJSON(w io.Writer, r []*webscraper.ScrapeResposne) error {
	e := json.NewEncoder(w)
	e.SetIndent("", "    ")
	return e.Encode(r)
}
