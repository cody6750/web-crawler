package webcrawler

import (
	"testing"

	webscraper "github.com/cody6750/codywebapi/webCrawler/webScraper"
)

func TestWebCrawler_Crawl(t *testing.T) {
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
		// 	name: "Crawl Correctly",
		// 	w:    crawl,
		// 	args: args{
		// 		url: "https://www.amazon.com/s?k=RTX+3080&ref=nb_sb_noss_2",
		// 		ScrapeURLConfiguration: []ScrapeURLConfiguration{
		// 			{
		// 				ExtractFromHTMLConfiguration: ExtractFromHTMLConfiguration{
		// 					Attribute:      "class",
		// 					AttributeValue: "a-link-normal",
		// 					Tag:            "a",
		// 				},
		// 				FormatURLConfiguration: FormatURLConfiguration{
		// 					PrefixToAdd: "http://amazon.com",
		// 				},
		// 			},
		// 		},
		// 	},
		// },
		{
			name: "Crawl Correctly",
			w:    crawl,
			args: args{
				url: "https://www.amazon.com/ZOTAC-Graphics-IceStorm-Advanced-ZT-A30800J-10PLHR/dp/B099ZCG8T5/ref=sr_1_4?keywords=RTX+3080&qid=1638491073&s=pc&sr=1-4",
				ScrapeURLConfiguration: []webscraper.ScrapeURLConfiguration{
					{
						// ExtractFromHTMLConfiguration: ExtractFromHTMLConfiguration{
						// 	Attribute:      "class",
						// 	AttributeValue: "a-link-normal",
						// 	Tag:            "a",
						// },
						FormatURLConfiguration: webscraper.FormatURLConfiguration{
							PrefixExist: "/",
							PrefixToAdd: "http://amazon.com",
						},
					},
				},
				ScrapeItemConfiguration: []webscraper.ScrapeItemConfiguration{
					{
						ItemName: "Product Name",
						ItemToGet: webscraper.ExtractFromHTMLConfiguration{
							Tag:            "div",
							Attribute:      "id",
							AttributeValue: "dp",
						},
						ItemDetails: map[string]webscraper.ExtractFromHTMLConfiguration{
							"price": {
								Tag:            "span",
								Attribute:      "aria-hidden",
								AttributeValue: "true",
							},
						},
					},
				},
			},
		},
		// {
		// 	name: "Crawl Correctly",
		// 	w:    crawl,
		// 	args: args{
		// 		url: "https://www.google.com/search?q=RTX+3080&rlz=1C1CHBF_enUS724US724&oq=RTX+3080&aqs=chrome.0.69i59j69i60l3j69i65.1086j0j15&sourceid=chrome&ie=UTF-8",
		// 		ScrapeURLConfiguration: []ScrapeURLConfiguration{
		// 			{
		// 				ExtractFromHTMLConfiguration: ExtractFromHTMLConfiguration{
		// 					Tag: "a",
		// 				},
		// 				FormatURLConfiguration: FormatURLConfiguration{
		// 					PrefixExist:    "/url?q=",
		// 					PrefixToRemove: "/url?q=",
		// 				},
		// 			},
		// 			{
		// 				ExtractFromHTMLConfiguration: ExtractFromHTMLConfiguration{
		// 					Tag: "a",
		// 				},
		// 				FormatURLConfiguration: FormatURLConfiguration{
		// 					PrefixExist: "/search",
		// 					PrefixToAdd: "https://www.google.com",
		// 				},
		// 			},
		// 		},
		// 	},
		// },
		// {
		// 	name: "Crawl Correctly",
		// 	w:    crawl,
		// 	args: args{
		// 		url: "https://www.newegg.com/p/pl?d=RTX+3080",
		// 	},
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.w.Crawl(tt.args.url, 1, tt.args.ScrapeItemConfiguration, tt.args.ScrapeURLConfiguration...)
		})
	}
}
