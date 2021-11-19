package webcrawler

import (
	"testing"

	webscraper "github.com/cody6750/codywebapi/webCrawler/webScraper"
)

func TestWebCrawler_Crawl(t *testing.T) {
	crawl := New()
	type args struct {
		url                 string
		ScrapeConfiguration []webscraper.ScrapeConfiguration
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
		// 		ScrapeConfiguration: []webscraper.ScrapeConfiguration{
		// 			{
		// 				ExtractURLFromHTMLConfiguration: webscraper.ExtractURLFromHTMLConfiguration{
		// 					AttributeToCheck:      "class",
		// 					AttributeValueToCheck: "a-link-normal",
		// 					TagToCheck:            "a",
		// 				},
		// 				FormatURLConfiguration: webscraper.FormatURLConfiguration{
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
				url: "https://www.google.com/search?q=RTX+3080&rlz=1C1CHBF_enUS724US724&oq=RTX+3080&aqs=chrome.0.69i59j69i60l3j69i65.1086j0j15&sourceid=chrome&ie=UTF-8",
				ScrapeConfiguration: []webscraper.ScrapeConfiguration{
					{
						ExtractURLFromHTMLConfiguration: webscraper.ExtractURLFromHTMLConfiguration{
							TagToCheck: "a",
						},
						FormatURLConfiguration: webscraper.FormatURLConfiguration{
							PrefixExist:    "/url?q=",
							PrefixToRemove: "/url?q=",
						},
					},
					{
						ExtractURLFromHTMLConfiguration: webscraper.ExtractURLFromHTMLConfiguration{
							TagToCheck: "a",
						},
						FormatURLConfiguration: webscraper.FormatURLConfiguration{
							PrefixExist: "/search",
							PrefixToAdd: "https://www.google.com",
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.w.Crawl(tt.args.url, tt.args.ScrapeConfiguration...)
		})
	}
}
