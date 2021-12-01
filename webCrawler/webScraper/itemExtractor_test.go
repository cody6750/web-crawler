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
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractItemWithScrapItemConfiguration(tt.args.token, tt.args.url, tt.args.itemTagsToCheck, tt.args.scrapeItemConfiguration)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractItemWithScrapItemConfiguration() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExtractItemWithScrapItemConfiguration() = %v, want %v", got, tt.want)
			}
		})
	}
}
