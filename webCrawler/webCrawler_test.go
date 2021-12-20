package webcrawler

import (
	"testing"

	webscraper "github.com/cody6750/codywebapi/webCrawler/webScraper"
)

func TestWebCrawler_Crawl(t *testing.T) {
	urlToCrawl := make(chan Queue, 500)
	duplicateUrls := make(chan map[string]bool, 2)
	q := Queue{}
	m := make(map[string]bool)
	duplicateUrls <- m
	urlToCrawl <- q
	crawl := New()
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
		// {
		// 	name: "Crawl Correctly- Product Page",
		// 	w:    crawl,
		// 	args: args{
		// 		url: "https://www.amazon.com/ZOTAC-Graphics-IceStorm-Advanced-ZT-A30800J-10PLHR/dp/B099ZCG8T5/ref=sr_1_4?keywords=RTX+3080&qid=1638491073&s=pc&sr=1-4",
		// 		ScrapeURLConfiguration: []webscraper.ScrapeURLConfiguration{
		// 			{
		// 				// ExtractFromHTMLConfiguration: ExtractFromHTMLConfiguration{
		// 				// 	Attribute:      "class",
		// 				// 	AttributeValue: "a-link-normal",
		// 				// 	Tag:            "a",
		// 				// },
		// 				FormatURLConfiguration: webscraper.FormatURLConfiguration{
		// 					PrefixExist: "/",
		// 					PrefixToAdd: "http://amazon.com",
		// 				},
		// 			},
		// 		},
		// 		ScrapeItemConfiguration: []webscraper.ScrapeItemConfiguration{
		// 			{
		// 				ItemName: "Product Name",
		// 				ItemToGet: webscraper.ExtractFromHTMLConfiguration{
		// 					Tag:            "div",
		// 					Attribute:      "id",
		// 					AttributeValue: "dp",
		// 				},
		// 				ItemDetails: map[string]webscraper.ExtractFromHTMLConfiguration{
		// 					"title": {
		// 						Tag:            "span",
		// 						Attribute:      "id",
		// 						AttributeValue: "productTitle",
		// 					},
		// 					"add-to-cart": {
		// 						Tag:            "span",
		// 						Attribute:      "id",
		// 						AttributeValue: "submit.add-to-cart-announce",
		// 					},
		// 					"ratings": {
		// 						Tag:            "span",
		// 						Attribute:      "id",
		// 						AttributeValue: "acrCustomerReviewText",
		// 					},
		// 					"price": {
		// 						Tag:            "span",
		// 						Attribute:      "class",
		// 						AttributeValue: "a-price a-text-price a-size-medium apexPriceToPay",
		// 					},
		// 				},
		// 			},
		// 		},
		// 	},
		// },
		// {
		// 	name: "Crawl Correctly - Product search page",
		// 	w:    crawl,
		// 	args: args{
		// 		url: "https://www.amazon.com/s?k=RTX+3080&ref=nb_sb_noss_2",
		// 		ScrapeURLConfiguration: []webscraper.ScrapeURLConfiguration{
		// 			{
		// 				// ExtractFromHTMLConfiguration: ExtractFromHTMLConfiguration{
		// 				// 	Attribute:      "class",
		// 				// 	AttributeValue: "a-link-normal",
		// 				// 	Tag:            "a",
		// 				// },
		// 				FormatURLConfiguration: webscraper.FormatURLConfiguration{
		// 					PrefixExist: "/",
		// 					PrefixToAdd: "http://amazon.com",
		// 				},
		// 			},
		// 		},
		// 		ScrapeItemConfiguration: []webscraper.ScrapeItemConfiguration{
		// 			{
		// 				ItemName: "Product Name",
		// 				ItemToGet: webscraper.ExtractFromHTMLConfiguration{
		// 					Tag:            "div",
		// 					Attribute:      "data-component-type",
		// 					AttributeValue: "s-search-result",
		// 				},
		// 				ItemDetails: map[string]webscraper.ExtractFromHTMLConfiguration{
		// 					"title": {
		// 						Tag:            "span",
		// 						Attribute:      "class",
		// 						AttributeValue: "a-size-medium a-color-base a-text-normal",
		// 					},
		// 					"price": {
		// 						Tag:            "span",
		// 						Attribute:      "class",
		// 						AttributeValue: "a-price",
		// 					},
		// 					"ratings": {
		// 						Tag:            "i",
		// 						Attribute:      "class",
		// 						AttributeValue: "a-icon a-icon-star-small a-star-small-4-5 aok-align-bottom",
		// 					},
		// 					"number of ratings": {
		// 						Tag:            "a",
		// 						Attribute:      "class",
		// 						AttributeValue: "a-link-normal",
		// 					},
		// 					"details": {
		// 						Tag:            "td",
		// 						Attribute:      "class",
		// 						AttributeValue: "a-size-base prodDetAttrValue",
		// 					},
		// 				},
		// 			},
		// 		},
		// 	},
		// },
		// {
		// 	name: "New egg Crawl Correctly - Product search page",
		// 	w:    crawl,
		// 	args: args{
		// 		url: "https://www.newegg.com/p/pl?d=rtx+3080&LeftPriceRange=1000+",
		// 		ScrapeURLConfiguration: []webscraper.ScrapeURLConfiguration{
		// 			{
		// 				// ExtractFromHTMLConfiguration: ExtractFromHTMLConfiguration{
		// 				// 	Attribute:      "class",
		// 				// 	AttributeValue: "a-link-normal",
		// 				// 	Tag:            "a",
		// 				// },
		// 				FormatURLConfiguration: webscraper.FormatURLConfiguration{
		// 					PrefixExist: "/",
		// 					PrefixToAdd: "http://newegg.com",
		// 				},
		// 			},
		// 		},
		// 		ScrapeItemConfiguration: []webscraper.ScrapeItemConfiguration{
		// 			{
		// 				ItemName: "Graphics Cards",
		// 				ItemToGet: webscraper.ExtractFromHTMLConfiguration{
		// 					Tag:            "div",
		// 					Attribute:      "class",
		// 					AttributeValue: "item-cell",
		// 				},
		// 				ItemDetails: map[string]webscraper.ExtractFromHTMLConfiguration{
		// 					"title": {
		// 						Tag:            "a",
		// 						Attribute:      "class",
		// 						AttributeValue: "item-title",
		// 					},
		// 					"price": {
		// 						Tag:            "strong",
		// 						Attribute:      "",
		// 						AttributeValue: "",
		// 					},
		// 					"link": {
		// 						Tag:            "a",
		// 						Attribute:      "class",
		// 						AttributeValue: "item-img",
		// 						AttributeToGet: "href",
		// 					},
		// 					"outofstock": {
		// 						Tag:            "i",
		// 						Attribute:      "class",
		// 						AttributeValue: "item-promo-icon",
		// 					},
		// 				},
		// 			},
		// 		},
		// 	},
		// },
		{
			name: "Best Buy Crawl Correctly - Product search page",
			w:    crawl,
			args: args{
				url: "https://www.bestbuy.com/site/searchpage.jsp?st=RTX+3080&_dyncharset=UTF-8&_dynSessConf=&id=pcat17071&type=page&sc=Global&cp=1&nrp=&sp=&qp=&list=n&af=true&iht=y&usc=All+Categories&ks=960&keys=keys",
				ScrapeURLConfiguration: []webscraper.ScrapeURLConfiguration{
					{
						// ExtractFromHTMLConfiguration: ExtractFromHTMLConfiguration{
						// 	Attribute:      "class",
						// 	AttributeValue: "a-link-normal",
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.w.Crawl(tt.args.url, duplicateUrls, 1, tt.args.ScrapeItemConfiguration, tt.args.ScrapeURLConfiguration...)
		})
	}
}
