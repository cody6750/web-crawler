package webcrawler

import (
	"testing"

	webscraper "github.com/cody6750/web-crawler/pkg/webScraper"
)

func TestWebCrawler_Crawl(t *testing.T) {
	crawl := NewCrawler()
	type args struct {
		url                     string
		ScrapeURLConfiguration  []webscraper.ScrapeURLConfig
		ScrapeItemConfiguration []webscraper.ScrapeItemConfig
	}
	tests := []struct {
		name string
		w    *WebCrawler
		args args
	}{
		// {
		// 	name: "Best Buy Crawl Correctly - Product search page",
		// 	w:    crawl,
		// 	args: args{
		// 		// url: "https://www.bestbuy.com/site/searchpage.jsp?id=pcat17071&qp=gpusv_facet%3DGraphics%20Processing%20Unit%20(GPU)~NVIDIA%20GeForce%20RTX%203080&st=rtx+3080",
		// 		url: "https://www.bestbuy.com/site/promo/pc-gaming-deals?qp=category_facet%3DComputer%20Cards%20%26%20Components~abcat0507000&sp=-bestsellingsort%20skuidsaas",
		// 		ScrapeURLConfiguration: []webscraper.ScrapeURLConfig{
		// 			{
		// 				// ExtractFromHTMLConfiguration: webscraper.ExtractFromHTMLConfiguration{
		// 				// 	Attribute:      "class",
		// 				// 	AttributeValue: "bottom-left-links",
		// 				// 	Tag:            "a",
		// 				// },
		// 				FormatURLConfiguration: webscraper.FormatURLConfiguration{
		// 					PrefixExist: "/",
		// 					PrefixToAdd: "http://bestbuy.com",
		// 				},
		// 			},
		// 		},
		// 		ScrapeItemConfiguration: []webscraper.ScrapeItemConfig{
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
		{
			name: "New Egg Crawl Correctly - Product search page",
			w:    crawl,
			args: args{
				url: "https://www.newegg.com/p/pl?d=RTX+3080",
				ScrapeURLConfiguration: []webscraper.ScrapeURLConfig{
					{
						FormatURLConfig: webscraper.FormatURLConfig{
							PrefixExist:    "////",
							PrefixToRemove: "////",
							PrefixToAdd:    "http://",
						},
					},
					{
						FormatURLConfig: webscraper.FormatURLConfig{
							PrefixExist:    "///",
							PrefixToRemove: "///",
							PrefixToAdd:    "http://",
						},
					},
					{
						FormatURLConfig: webscraper.FormatURLConfig{
							PrefixExist:    "//",
							PrefixToRemove: "//",
							PrefixToAdd:    "http://",
						},
					},
					{
						FormatURLConfig: webscraper.FormatURLConfig{
							PrefixExist: "/",
							PrefixToAdd: "http://newegg.com",
						},
					},
				},
				ScrapeItemConfiguration: []webscraper.ScrapeItemConfig{
					{
						ItemName: "Graphics Cards",
						ItemToGet: webscraper.ExtractFromTokenConfig{
							Tag:            "div",
							Attribute:      "class",
							AttributeValue: "item-container",
						},
						ItemDetails: map[string]webscraper.ExtractFromTokenConfig{
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
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.w.Crawl(tt.args.url, tt.args.ScrapeItemConfiguration, tt.args.ScrapeURLConfiguration...)
		})
	}
}
