package data

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	webcrawler "github.com/cody6750/web-crawler/pkg"
	webscraper "github.com/cody6750/web-crawler/pkg/webScraper"
)

// Product ...
type Product struct {
	Name                         string            `json:"name"`
	Description                  string            `json:"description"`
	Price                        float32           `json:"price"`
	URL                          string            `json:"url"`
	ParentURL                    string            `json:"-"`
	AdditionalProductInformation map[string]string `json:"additional_Product_Information"`
	CreatedOn                    string            `json:"-"`
	UpdatedOn                    string            `json:"-"`
	DeletedOn                    string            `json:"-"`
}

// ProductPayload ...
type ProductPayload struct {
	ScrapeItemConfiguration []webscraper.ScrapeItemConfiguration `json:"ScrapeItemConfiguration"`
	ScrapeURLConfiguration  []webscraper.ScrapeURLConfiguration  `json:"ScrapeURLConfiguration"`
	RootURL                 string                               `json:"RootURL"`
}

// Products ...
type Products []*Product

// MarshalPayloadToStruct ...
func MarshalPayloadToStruct(rw http.ResponseWriter, r *http.Request) (*ProductPayload, error) {
	payload := &ProductPayload{}
	err := json.NewDecoder(r.Body).Decode(payload)
	if err != nil {
		log.Print(err)
		return payload, err
	}
	return payload, nil
}

// GetProduct ...
func GetProduct(url string, itemsToget []webscraper.ScrapeItemConfiguration, ScrapeURLConfiguration ...webscraper.ScrapeURLConfiguration) ([]*webscraper.ScrapeResposne, error) {
	crawler := webcrawler.NewCrawler()
	log.Print(itemsToget)
	response, err := crawler.Crawl(url, itemsToget, ScrapeURLConfiguration...)
	// response, err := crawler.Crawl("https://www.newegg.com/p/pl?d=RTX+3080",
	// 	[]webscraper.ScrapeItemConfiguration{
	// 		{
	// 			ItemName: "Graphics Cards",
	// 			ItemToGet: webscraper.ExtractFromHTMLConfiguration{
	// 				Tag:            "div",
	// 				Attribute:      "class",
	// 				AttributeValue: "item-container",
	// 			},
	// 			ItemDetails: map[string]webscraper.ExtractFromHTMLConfiguration{
	// 				"title": {
	// 					Tag:            "a",
	// 					Attribute:      "class",
	// 					AttributeValue: "item-title",
	// 				},
	// 				"price": {
	// 					Tag:            "strong",
	// 					Attribute:      "",
	// 					AttributeValue: "",
	// 				},
	// 			},
	// 		},
	// 	},
	// 	[]webscraper.ScrapeURLConfiguration{
	// 		{
	// 			FormatURLConfiguration: webscraper.FormatURLConfiguration{
	// 				PrefixExist:    "////",
	// 				PrefixToRemove: "////",
	// 				PrefixToAdd:    "http://",
	// 			},
	// 		},
	// 		{
	// 			FormatURLConfiguration: webscraper.FormatURLConfiguration{
	// 				PrefixExist:    "///",
	// 				PrefixToRemove: "///",
	// 				PrefixToAdd:    "http://",
	// 			},
	// 		},
	// 		{
	// 			FormatURLConfiguration: webscraper.FormatURLConfiguration{
	// 				PrefixExist:    "//",
	// 				PrefixToRemove: "//",
	// 				PrefixToAdd:    "http://",
	// 			},
	// 		},
	// 		{
	// 			FormatURLConfiguration: webscraper.FormatURLConfiguration{
	// 				PrefixExist: "/",
	// 				PrefixToAdd: "http://newegg.com",
	// 			},
	// 		},
	// 	}...,
	// )
	return response, err
}

// ToJSON ...
func ToJSON(w io.Writer, r []*webscraper.ScrapeResposne) error {
	// bytes, _ := json.MarshalIndent(r, "", "    ")
	// log.Print(string(bytes))
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "    ")
	return encoder.Encode(r)
}
