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
		// 					AttributeToCheck:      "class",
		// 					AttributeValueToCheck: "a-link-normal",
		// 					TagToCheck:            "a",
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
				url: "https://www.amazon.com/gp/offer-listing/B08W8DGK3X/ref=dp_olp_unknown_mbc",
				ScrapeURLConfiguration: []webscraper.ScrapeURLConfiguration{
					{
						// ExtractFromHTMLConfiguration: ExtractFromHTMLConfiguration{
						// 	AttributeToCheck:      "class",
						// 	AttributeValueToCheck: "a-link-normal",
						// 	TagToCheck:            "a",
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
						ItemToget: map[string]webscraper.ExtractFromHTMLConfiguration{
							"price": webscraper.ExtractFromHTMLConfiguration{
								TagToCheck:            "span",
								AttributeToCheck:      "class",
								AttributeValueToCheck: "a-price-whole",
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
		// 					TagToCheck: "a",
		// 				},
		// 				FormatURLConfiguration: FormatURLConfiguration{
		// 					PrefixExist:    "/url?q=",
		// 					PrefixToRemove: "/url?q=",
		// 				},
		// 			},
		// 			{
		// 				ExtractFromHTMLConfiguration: ExtractFromHTMLConfiguration{
		// 					TagToCheck: "a",
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
