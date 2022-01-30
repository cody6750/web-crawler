package data

import (
	"encoding/json"
	"io"
	"log"

	webcrawler "github.com/cody6750/codywebapi/webCrawler"
	webscraper "github.com/cody6750/codywebapi/webCrawler/webScraper"
)

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

type Products []*Product

func GetProduct() ([]*webscraper.ScrapeResposne, error) {
	crawler := webcrawler.NewCrawler()
	response, err := crawler.Crawl("https://www.newegg.com/p/pl?d=RTX+3080",
		[]webscraper.ScrapeItemConfiguration{
			{
				ItemName: "Graphics Cards",
				ItemToGet: webscraper.ExtractFromHTMLConfiguration{
					Tag:            "div",
					Attribute:      "class",
					AttributeValue: "item-container",
				},
				ItemDetails: map[string]webscraper.ExtractFromHTMLConfiguration{
					"title": {
						Tag:            "a",
						Attribute:      "class",
						AttributeValue: "item-title",
					},
					"price": {
						Tag:            "strong",
						Attribute:      "",
						AttributeValue: "",
					},
				},
			},
		},
		[]webscraper.ScrapeURLConfiguration{
			{
				FormatURLConfiguration: webscraper.FormatURLConfiguration{
					PrefixExist:    "////",
					PrefixToRemove: "////",
					PrefixToAdd:    "http://",
				},
			},
			{
				FormatURLConfiguration: webscraper.FormatURLConfiguration{
					PrefixExist:    "///",
					PrefixToRemove: "///",
					PrefixToAdd:    "http://",
				},
			},
			{
				FormatURLConfiguration: webscraper.FormatURLConfiguration{
					PrefixExist:    "//",
					PrefixToRemove: "//",
					PrefixToAdd:    "http://",
				},
			},
			{
				FormatURLConfiguration: webscraper.FormatURLConfiguration{
					PrefixExist: "/",
					PrefixToAdd: "http://newegg.com",
				},
			},
		}...,
	)
	return response, err
}

func ToJSON(w io.Writer, r []*webscraper.ScrapeResposne) error {
	bytes, _ := json.MarshalIndent(r, "", "    ")
	log.Print(string(bytes))
	return json.NewEncoder(w).Encode(r)
}
