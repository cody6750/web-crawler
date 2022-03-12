package webcrawler

import (
	"log"
	"testing"

	"golang.org/x/net/html"
)

func Test_extractURLFromHTMLUsingConfiguration(t *testing.T) {
	type args struct {
		token     html.Token
		urlConfig ExtractFromTokenConfig
	}
	tests := []struct {
		name               string
		args               args
		wantAttributeValue string
		wantErr            error
	}{
		{
			name: "Extract URL from HTML Using Configuration",
			args: args{
				token: html.Token{
					Data: "span",
					Attr: []html.Attribute{
						{Key: "class", Val: "a-link-normal"},
						{Key: "href", Val: "amazon.com"},
					},
				},
				urlConfig: ExtractFromTokenConfig{
					Tag:            "span",
					Attribute:      "class",
					AttributeValue: "a-link-normal",
				},
			},
			wantAttributeValue: "amazon.com",
			wantErr:            nil,
		},
		{
			name: "Fail to Extract URL from HTML Using Configuration",
			args: args{
				token: html.Token{
					Data: "span",
					Attr: []html.Attribute{
						{Key: "class", Val: "a-link-normal"},
						{Key: "href", Val: "amazon.com"},
					},
				},
				urlConfig: ExtractFromTokenConfig{
					Tag:            "",
					Attribute:      "class",
					AttributeValue: "not it sis",
				},
			},
			wantAttributeValue: "",
			wantErr:            nil,
		},
	}

	for _, tt := range tests {
		log.Printf("[TEST]: %v has started\n", tt.name)
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractURLFromTokenUsingConfig(tt.args.token, tt.args.urlConfig)
			if got != tt.wantAttributeValue {
				log.Printf("[TEST]: %v has failed\n\n", tt.name)
				t.Errorf("extractURLFromTokenUsingConfig() = %v, want %v", got, tt.wantAttributeValue)
				return
			}
			if err != tt.wantErr {
				log.Printf("[TEST]: %v has failed\n\n", tt.name)
				t.Errorf("extractURLFromTokenUsingConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			log.Printf("[TEST]: %v has successfully finished\n\n", tt.name)

		})
	}
}

func Test_extractURLFromHTML(t *testing.T) {
	type args struct {
		token html.Token
	}
	tests := []struct {
		name    string
		args    args
		wantURL string
		wantErr error
	}{
		{
			name: "Get HTTP Attribute",
			args: args{
				token: html.Token{

					Attr: []html.Attribute{
						{Key: "class", Val: "a-link-normal"},
						{Key: "href", Val: "amazon.com"},
					},
				},
			},
			wantURL: "amazon.com",
			wantErr: nil,
		},
		{
			name: "Get Wrong HTTP Attribute",
			args: args{
				token: html.Token{
					Attr: []html.Attribute{
						{Key: "class", Val: "a-link-normals"},
					},
				},
			},
			wantURL: "",
			wantErr: nil,
		}}
	for _, tt := range tests {
		log.Printf("[TEST]: %v has started\n", tt.name)
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractURLFromToken(tt.args.token)
			if got != tt.wantURL {
				log.Printf("[TEST]: %v has failed\n\n", tt.name)
				t.Errorf("extractURLFromHTML() = %v, want %v", got, tt.wantURL)
			}
			if err != tt.wantErr {
				t.Errorf("extractURLFromHTML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			log.Printf("extractURLFromHTML() = %v, want %v", got, tt.wantURL)
			log.Printf("[TEST]: %v has successfully finished\n\n", tt.name)

		})
	}
}

func Test_getHTTPAttributeValueFromToken(t *testing.T) {
	type args struct {
		token          html.Token
		attributeToGet string
	}
	tests := []struct {
		name               string
		args               args
		wantAttributeValue string
		wantErr            error
	}{
		{
			name: "Get HTTP Attribute",
			args: args{
				token: html.Token{
					Attr: []html.Attribute{
						{Key: "class", Val: "a-link-normal"},
					},
				},
				attributeToGet: "class",
			},
			wantAttributeValue: "a-link-normal",
			wantErr:            nil,
		},
		{
			name: "Get Wrong HTTP Attribute",
			args: args{
				token: html.Token{
					Attr: []html.Attribute{
						{Key: "class", Val: "a-link-normals"},
					},
				},
				attributeToGet: "class",
			},
			wantAttributeValue: "a-link-normal",
			wantErr:            nil,
		},
	}
	for _, tt := range tests {
		log.Printf("[TEST]: %v has started\n", tt.name)
		t.Run(tt.name, func(t *testing.T) {
			gotAttributeValue, err := extractAttributeValue(tt.args.token, tt.args.attributeToGet)
			if gotAttributeValue != tt.wantAttributeValue {
				log.Printf("[TEST]: %v has failed\n\n", tt.name)
				t.Errorf("getHTTPAttributeValueFromToken() = %v, want %v", gotAttributeValue, tt.wantAttributeValue)
				return
			}
			if err != tt.wantErr {
				t.Errorf("getHTTPAttributeValueFromToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			log.Printf("isEmptyScrapeURLConfiguration() = %v, want %v", gotAttributeValue, tt.wantAttributeValue)
			log.Printf("[TEST]: %v has successfully finished\n\n", tt.name)
		})
	}
}

func Test_formatURL(t *testing.T) {
	type args struct {
		url             string
		formatURLConfig FormatURLConfig
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr error
	}{
		{
			name: "Formated correctly",
			args: args{
				url: "ReplacePrefix/RTX_3070/ReplaceSuffix",
				formatURLConfig: FormatURLConfig{
					SuffixToAdd:      "checkout",
					SuffixToRemove:   "ReplaceSuffix",
					PrefixToAdd:      "amazon.com",
					PrefixToRemove:   "ReplacePrefix",
					ReplaceOldString: "RTX_3070",
					ReplaceNewString: "RTX_3080",
				},
			},
			want:    "amazon.com/RTX_3080/checkout",
			wantErr: nil,
		},
		{
			name: "Formated incorrectly",
			args: args{
				url: "ReplaceSuffix/RTX_3070/ReplacePrefix",
				formatURLConfig: FormatURLConfig{
					SuffixToAdd:      "DOESNOTWORK",
					SuffixToRemove:   "DOESNOTWORK",
					PrefixToAdd:      "DOESNOTWORK",
					PrefixToRemove:   "DOESNOTWORK",
					ReplaceOldString: "DOESNOTWORK",
					ReplaceNewString: "DOESNOTWORK",
				},
			},
			want:    "",
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		log.Printf("[TEST]: %v has started\n", tt.name)
		t.Run(tt.name, func(t *testing.T) {
			got := formatURL(tt.args.url, tt.args.formatURLConfig)
			if got != tt.want {
				log.Printf("[TEST]: %v has failed\n\n", tt.name)
				t.Errorf("formatURL() = %v, want %v", got, tt.want)
				return
			}
			log.Printf("isEmptyScrapeURLConfiguration() = %v, want %v", got, tt.want)
			log.Printf("[TEST]: %v has successfully finished\n\n", tt.name)
		})
	}
}

func Test_isEmptyFormatURLConfiguration(t *testing.T) {
	type args struct {
		formatURLConfiguration FormatURLConfig
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Is empty",
			want: true,
			args: args{
				formatURLConfiguration: FormatURLConfig{},
			},
		},
		{
			name: "Is not empty",
			want: false,
			args: args{
				formatURLConfiguration: FormatURLConfig{
					SuffixToRemove: "not empty",
				},
			},
		},
	}
	for _, tt := range tests {
		log.Printf("[TEST]: %v has started\n", tt.name)
		t.Run(tt.name, func(t *testing.T) {
			if got := IsEmpty(tt.args.formatURLConfiguration); got != tt.want {
				log.Printf("[TEST]: %v has failed\n\n", tt.name)
				t.Errorf("isEmptyFormatURLConfiguration() = %v, want %v", got, tt.want)
			} else {
				log.Printf("isEmptyScrapeURLConfiguration() = %v, want %v", got, tt.want)
				log.Printf("[TEST]: %v has successfully finished\n\n", tt.name)
			}
		})
	}
}

func Test_isEmptyExtractFromHTMLConfiguration(t *testing.T) {
	type args struct {
		extractFromHTMLConfiguration ExtractFromTokenConfig
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Is empty",
			args: args{
				extractFromHTMLConfiguration: ExtractFromTokenConfig{},
			},
			want: true,
		},
		{
			name: "Is not empty",
			args: args{
				extractFromHTMLConfiguration: ExtractFromTokenConfig{
					Attribute: "hi",
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		log.Printf("[TEST]: %v has started\n", tt.name)
		t.Run(tt.name, func(t *testing.T) {
			if got := IsEmpty(tt.args.extractFromHTMLConfiguration); got != tt.want {
				log.Printf("[TEST]: %v has failed\n\n", tt.name)
				t.Errorf("isEmptyExtractFromHTMLConfiguration() = test name:%v ,%v, want %v", tt.name, got, tt.want)
			} else {
				log.Printf("isEmptyScrapeURLConfiguration() = %v, want %v", got, tt.want)
				log.Printf("[TEST]: %v has successfully finished\n\n", tt.name)
			}
		})
	}
}

func Test_isEmptyScrapeURLConfiguration(t *testing.T) {
	type args struct {
		s []ScrapeURLConfig
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Not empty scrap configuration slice",
			args: args{
				s: []ScrapeURLConfig{
					{
						Name:                   "",
						ExtractFromTokenConfig: ExtractFromTokenConfig{},
						FormatURLConfig:        FormatURLConfig{},
					},
				},
			},
			want: false,
		},
		{
			name: "Empty scrap configuration slice",
			args: args{},
			want: true,
		},
	}
	for _, tt := range tests {
		log.Printf("[TEST]: %v has started\n", tt.name)
		t.Run(tt.name, func(t *testing.T) {
			if got := IsEmpty(tt.args.s); got != tt.want {
				log.Printf("[TEST]: %v has failed\n\n", tt.name)
				t.Errorf("isEmptyScrapeURLConfiguration() = %v, want %v", got, tt.want)
			} else {
				log.Printf("isEmptyScrapeURLConfiguration() = %v, want %v", got, tt.want)
				log.Printf("[TEST]: %v has successfully finished\n\n", tt.name)
			}
		})
	}
}

func Test_isDuplicateURL(t *testing.T) {
	type args struct {
		url         string
		URLsToCheck map[string]bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Has Duplicate Url",
			args: args{
				url: "exist",
				URLsToCheck: map[string]bool{
					"exist": true,
				},
			},
			want: true,
		},
		{
			name: "Has No Duplicate Url",
			args: args{
				url:         "does not exist",
				URLsToCheck: map[string]bool{},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		log.Printf("[TEST]: %v has started\n", tt.name)
		t.Run(tt.name, func(t *testing.T) {
			if got := isDuplicateURL(tt.args.url, tt.args.URLsToCheck); got != tt.want {
				log.Printf("[TEST]: %v has failed\n\n", tt.name)
				t.Errorf("isDuplicateURL() = %v, want %v", got, tt.want)
			} else {
				log.Printf("isDuplicateURL() = %v, want %v", got, tt.want)
				log.Printf("[TEST]: %v has successfully finished\n\n", tt.name)
			}
		})
	}
}

func TestWebScraper_checkBlackListedURLPaths(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		w    *WebScraper
		args args
		want bool
	}{
		{
			name: "blacklisted",
			w: &WebScraper{
				BlackListedURLPaths: map[string]struct{}{
					"/gp/cart": {},
				},
			},
			args: args{
				url: "https://www.amazon.com/gp/cart/view.html?ref_=nav_cart",
			},
			want: true,
		}, {
			name: "Not blacklistd",
			w: &WebScraper{
				BlackListedURLPaths: map[string]struct{}{
					"/gp/cart": {},
				},
			},
			args: args{
				url: "https://www.amazon.com/gp/view.html?ref_=nav_cart",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.w.isBlackListedURLPath(tt.args.url); got != tt.want {
				t.Errorf("WebScraper.checkBlackListedURLPaths() = %v, want %v", got, tt.want)
			}
		})
	}
}
