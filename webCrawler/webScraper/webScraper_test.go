package webcrawler

import (
	"log"
	"testing"

	"golang.org/x/net/html"
)

func Test_extractURLFromHTMLUsingConfiguration(t *testing.T) {
	type args struct {
		token     html.Token
		urlConfig ExtractFromHTMLConfiguration
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
				urlConfig: ExtractFromHTMLConfiguration{
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
				urlConfig: ExtractFromHTMLConfiguration{
					Tag:            "",
					Attribute:      "class",
					AttributeValue: "not it sis",
				},
			},
			wantAttributeValue: "",
			wantErr:            errExtractURLFromHTMLUsingConfiguration,
		},
	}

	for _, tt := range tests {
		log.Printf("[TEST]: %v has started\n", tt.name)
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractURLFromHTMLUsingConfiguration(tt.args.token, tt.args.urlConfig)
			if got != tt.wantAttributeValue {
				log.Printf("[TEST]: %v has failed\n\n", tt.name)
				t.Errorf("extractURLFromHTMLUsingConfiguration() = %v, want %v", got, tt.wantAttributeValue)
				return
			}
			if err != tt.wantErr {
				log.Printf("[TEST]: %v has failed\n\n", tt.name)
				t.Errorf("extractURLFromHTMLUsingConfiguration() error = %v, wantErr %v", err, tt.wantErr)
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
			wantErr: errExtractURLFromHTML,
		}}
	for _, tt := range tests {
		log.Printf("[TEST]: %v has started\n", tt.name)
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractURLFromHTML(tt.args.token)
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
			gotAttributeValue, err := getHTTPAttributeValueFromToken(tt.args.token, tt.args.attributeToGet)
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
		formatURLConfig FormatURLConfiguration
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
				formatURLConfig: FormatURLConfiguration{
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
				formatURLConfig: FormatURLConfiguration{
					SuffixToAdd:      "DOESNOTWORK",
					SuffixToRemove:   "DOESNOTWORK",
					PrefixToAdd:      "DOESNOTWORK",
					PrefixToRemove:   "DOESNOTWORK",
					ReplaceOldString: "DOESNOTWORK",
					ReplaceNewString: "DOESNOTWORK",
				},
			},
			want:    "",
			wantErr: errFormatURL,
		},
	}
	for _, tt := range tests {
		log.Printf("[TEST]: %v has started\n", tt.name)
		t.Run(tt.name, func(t *testing.T) {
			got, err := formatURL(tt.args.url, tt.args.formatURLConfig)
			if err != tt.wantErr {
				log.Printf("[TEST]: %v has failed\n\n", tt.name)
				t.Errorf("formatURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
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
		formatURLConfiguration FormatURLConfiguration
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
				formatURLConfiguration: FormatURLConfiguration{},
			},
		},
		{
			name: "Is not empty",
			want: false,
			args: args{
				formatURLConfiguration: FormatURLConfiguration{
					SuffixToRemove: "not empty",
				},
			},
		},
	}
	for _, tt := range tests {
		log.Printf("[TEST]: %v has started\n", tt.name)
		t.Run(tt.name, func(t *testing.T) {
			if got := isEmptyFormatURLConfiguration(tt.args.formatURLConfiguration); got != tt.want {
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
		extractFromHTMLConfiguration ExtractFromHTMLConfiguration
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Is empty",
			args: args{
				extractFromHTMLConfiguration: ExtractFromHTMLConfiguration{},
			},
			want: true,
		},
		{
			name: "Is not empty",
			args: args{
				extractFromHTMLConfiguration: ExtractFromHTMLConfiguration{
					Attribute: "hi",
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		log.Printf("[TEST]: %v has started\n", tt.name)
		t.Run(tt.name, func(t *testing.T) {
			if got := isEmptyExtractFromHTMLConfiguration(tt.args.extractFromHTMLConfiguration); got != tt.want {
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
		s []ScrapeURLConfiguration
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Not empty scrap configuration slice",
			args: args{
				s: []ScrapeURLConfiguration{
					{
						ConfigurationName:            "",
						ExtractFromHTMLConfiguration: ExtractFromHTMLConfiguration{},
						FormatURLConfiguration:       FormatURLConfiguration{},
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
			if got := isEmptyScrapeURLConfiguration(tt.args.s); got != tt.want {
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
