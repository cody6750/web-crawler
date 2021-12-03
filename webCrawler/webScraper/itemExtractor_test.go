package webcrawler

import (
	"reflect"
	"testing"

	"golang.org/x/net/html"
)

func TestExtractItemWithScrapItemConfiguration(t *testing.T) {
	type args struct {
		token                   html.Token
		url                     string
		itemTagsToCheck         map[string]bool
		scrapeItemConfiguration []ScrapeItemConfiguration
	}
	tests := []struct {
		name    string
		args    args
		want    Item
		wantErr error
	}{
		{
			name: "Generate Item",
			args: args{
				token: html.Token{
					Data: "span",
					Attr: []html.Attribute{
						{
							Key: "class",
							Val: "a-price-whole",
						},
					},
				},
				url: "https://www.amazon.com/gp/offer-listing/B08W8DGK3X/ref=dp_olp_unknown_mbc",
				itemTagsToCheck: map[string]bool{
					"span": true,
				},
				scrapeItemConfiguration: []ScrapeItemConfiguration{
					{
						ItemName: "RTX 3080",
						URL:      "https://www.amazon.com/gp/offer-listing/B08W8DGK3X/ref=dp_olp_unknown_mbc",
						ItemToget: map[string]ExtractFromHTMLConfiguration{
							"price": {
								Tag:            "span",
								Attribute:      "class",
								AttributeValue: "a-price-whole",
							},
						},
					},
				},
			},
			want: Item{
				ItemName:    "RTX 3080",
				URL:         "https://www.amazon.com/gp/offer-listing/B08W8DGK3X/ref=dp_olp_unknown_mbc",
				ItemDetails: "",
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractItemWithScrapItemConfiguration(tt.args.token, tt.args.url, tt.args.itemTagsToCheck, tt.args.scrapeItemConfiguration)
			if err != tt.wantErr {
				t.Errorf("ExtractItemWithScrapItemConfiguration() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExtractItemWithScrapItemConfiguration() = %v, want %v", got, tt.want)
			}
		})
	}
}
