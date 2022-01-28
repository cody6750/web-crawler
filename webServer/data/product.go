package data

import (
	"encoding/json"
	"io"

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

func GetProduct() Products {
	crawler := webcrawler.NewCrawler()
	crawler.Crawl("https://www.bestbuy.com/site/searchpage.jsp?id=pcat17071&qp=gpusv_facet%3DGraphics%20Processing%20Unit%20(GPU)~NVIDIA%20GeForce%20RTX%203080&st=rtx+3080",
		[]webscraper.ScrapeItemConfiguration{
			{
				ItemName: "Graphics Cards",
				ItemToGet: webscraper.ExtractFromHTMLConfiguration{
					Tag:            "li",
					Attribute:      "class",
					AttributeValue: "sku-item",
				},
				ItemDetails: map[string]webscraper.ExtractFromHTMLConfiguration{
					"title": {
						Tag:            "h4",
						Attribute:      "class",
						AttributeValue: "sku-header",
					},
					"price": {
						Tag:            "span",
						Attribute:      "aria-hidden",
						AttributeValue: "true",
					},
					"link": {
						Tag:            "a",
						Attribute:      "",
						AttributeValue: "",
						AttributeToGet: "href",
					},
					"In stock": {
						Tag:            "button",
						AttributeToGet: "data-button-state",
						AttributeValue: "button",
						Attribute:      "disabled type",
					},
					"Out of stock": {
						Tag:            "button",
						AttributeToGet: "data-button-state",
						AttributeValue: "button",
						Attribute:      "type",
					},
				},
			},
		},
		[]webscraper.ScrapeURLConfiguration{
			{
				FormatURLConfiguration: webscraper.FormatURLConfiguration{
					PrefixExist: "/",
					PrefixToAdd: "http://bestbuy.com",
				},
			},
		}...,
	)
	var products Products
	return products
}

func (p *Products) ToJSON(w io.Writer) error {
	return json.NewEncoder(w).Encode(p)
}
