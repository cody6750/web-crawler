package webcrawler

import (
	"testing"

	webscraper "github.com/cody6750/codywebapi/webCrawler/webScraper"
)

func TestWebCrawler_Crawl(t *testing.T) {
	crawl := NewCrawler()
	type args struct {
		url                     string
		ScrapeURLConfiguration  []webscraper.ScrapeURLConfiguration
		ScrapeItemConfiguration []webscraper.ScrapeItemConfiguration
	}
	tests := []struct {
		name string
		w    *WebCrawler
		args args
	}{
		{
			name: "Best Buy Crawl Correctly - Product search page",
			w:    crawl,
			args: args{
				url: "https://www.bestbuy.com/site/searchpage.jsp?id=pcat17071&qp=gpusv_facet%3DGraphics%20Processing%20Unit%20(GPU)~NVIDIA%20GeForce%20RTX%203080&st=rtx+3080",
				ScrapeURLConfiguration: []webscraper.ScrapeURLConfiguration{
					{
						// ExtractFromHTMLConfiguration: webscraper.ExtractFromHTMLConfiguration{
						// 	Attribute:      "class",
						// 	AttributeValue: "bottom-left-links",
						// 	Tag:            "a",
						// },
						FormatURLConfiguration: webscraper.FormatURLConfiguration{
							PrefixExist: "/",
							PrefixToAdd: "http://bestbuy.com",
						},
					},
				},
				ScrapeItemConfiguration: []webscraper.ScrapeItemConfiguration{
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
			},
		},
		// {
		// 	name: "Best Buy Crawl Correctly - Product search page",
		// 	w:    crawl,
		// 	args: args{
		// 		url: "https://www.bestbuy.com/site/searchpage.jsp?st=RTX+3080&_dyncharset=UTF-8&_dynSessConf=&id=pcat17071&type=page&sc=Global&cp=1&nrp=&sp=&qp=&list=n&af=true&iht=y&usc=All+Categories&ks=960&keys=keys",
		// 		ScrapeURLConfiguration: []webscraper.ScrapeURLConfiguration{
		// 			{
		// 				// ExtractFromHTMLConfiguration: ExtractFromHTMLConfiguration{
		// 				// 	Attribute:      "class",
		// 				// 	AttributeValue: "a-link-normal",
		// 				// 	Tag:            "a",
		// 				// },
		// 				FormatURLConfiguration: webscraper.FormatURLConfiguration{
		// 					PrefixExist: "/",
		// 					PrefixToAdd: "http://bestbuy.com",
		// 				},
		// 			},
		// 		},
		// 		ScrapeItemConfiguration: []webscraper.ScrapeItemConfiguration{
		// 			{
		// 				ItemName: "Graphics Cards",
		// 				ItemToGet: webscraper.ExtractFromHTMLConfiguration{
		// 					Tag:            "li",
		// 					Attribute:      "class",
		// 					AttributeValue: "sku-item",
		// 				},
		// 				ItemDetails: map[string]webscraper.ExtractFromHTMLConfiguration{
		// 					"title": {
		// 						Tag:            "h4",
		// 						Attribute:      "class",
		// 						AttributeValue: "sku-header",
		// 					},
		// 					"price": {
		// 						Tag:            "span",
		// 						Attribute:      "aria-hidden",
		// 						AttributeValue: "true",
		// 					},
		// 					"link": {
		// 						Tag:            "a",
		// 						Attribute:      "",
		// 						AttributeValue: "",
		// 						AttributeToGet: "href",
		// 					},
		// 					"In stock": {
		// 						Tag:            "button",
		// 						AttributeToGet: "data-button-state",
		// 						AttributeValue: "button",
		// 						Attribute:      "disabled type",
		// 					},
		// 					"Out of stock": {
		// 						Tag:            "button",
		// 						AttributeToGet: "data-button-state",
		// 						AttributeValue: "button",
		// 						Attribute:      "type",
		// 					},
		// 				},
		// 			},
		// 		},
		// 	},
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.w.Crawl(tt.args.url, tt.args.ScrapeItemConfiguration, tt.args.ScrapeURLConfiguration...)
		})
	}
}
