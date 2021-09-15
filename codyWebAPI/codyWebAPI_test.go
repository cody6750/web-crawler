package codywebapi

import (
	"log"
	"os"
	"testing"
)

const (
	lengthTest             string = "Length test"
	noFlags                string = "No flags test"
	unsupportedActionTest  string = "Unsupported action test"
	unsupportedFlagTest    string = "Unsupported flag test"
	unsupportedParseFlags  string = "Unsupported parse flag test"
	unsupportedWebsiteTest string = "Unsupported website test"
	supportedSearchTest    string = "Supported search test"
)

func Test_parseInput(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		args     args
		expected error
		name     string
	}{
		{
			name:     lengthTest,
			expected: errInput,
			args: args{
				input: "codyWebAPI",
			},
		},
		{
			name:     unsupportedActionTest,
			expected: errInput,
			args: args{
				input: "codyWebAPI notsupported",
			},
		},
		{
			name:     unsupportedWebsiteTest,
			expected: errWebsiteFlag,
			args: args{
				input: "codyWebAPI search --website notsupported",
			},
		},
		{
			name:     unsupportedFlagTest,
			expected: errUnsupportedFlag,
			args: args{
				input: "codyWebAPI search --notsupported notsupported",
			},
		},
		{
			name:     supportedSearchTest,
			expected: nil,
			args: args{
				input: "codyWebAPI search --website amazon --item GTX 1080",
			},
		},
		{
			name:     "Success",
			expected: nil,
			args: args{
				input: "codyWebAPI search --item GTX 1080 --website amazon",
			},
		},
	}
	for _, tt := range tests {
		log.Printf("[TEST]: %v has started\n", tt.name)
		t.Run(tt.name, func(t *testing.T) {
			actual := parseInput(tt.args.input)
			if actual != tt.expected {
				log.Printf("[TEST]: %v has failed want: %v got: %v\n\n", tt.name, tt.expected, actual)
				t.Errorf("Failed test %v", tt.name)
			} else {
				log.Printf("[TEST]: %v has successfully finished\n\n", tt.name)
			}
		})
	}
}

func Test_parseFlags(t *testing.T) {
	type args struct {
		applicationCommand string
		input              []string
	}
	tests := []struct {
		expected      []string
		expectedError error
		name          string
		args          args
	}{
		{
			name:          noFlags,
			expected:      []string{"", ""},
			expectedError: errParseFlag,
			args: args{
				applicationCommand: "search",
				input:              []string{"amazon", "--i", "GTX", "1080"},
			},
		},
		{
			name:          "flagNotSupported",
			expected:      []string{"", ""},
			expectedError: errUnsupportedFlag,
			args: args{
				applicationCommand: "search",
				input:              []string{"--novalue1", "--novalue2"},
			},
		},
		{
			name:          "flagNotSet",
			expected:      []string{"", ""},
			expectedError: errFlagNotSet,
			args: args{
				applicationCommand: "search",
				input:              []string{"--website", "--novalue2"},
			},
		},
		{
			name:          unsupportedParseFlags,
			expected:      []string{"", ""},
			expectedError: errUnsupportedFlag,
			args: args{
				applicationCommand: "search",
				input:              []string{"--w", "amazon", "--i", "GTX", "1080"},
			},
		},
		{
			name:          unsupportedParseFlags,
			expected:      []string{"", ""},
			expectedError: errUnsupportedFlag,
			args: args{
				applicationCommand: "search",
				input:              []string{"--website", "amazon", "--i", "GTX", "1080"},
			},
		},
		{
			name:          "supportedParseFlags",
			expected:      []string{"amazon", "GTX 1080"},
			expectedError: nil,
			args: args{
				applicationCommand: "search",
				input:              []string{"--website", "amazon", "--item", "GTX", "1080"},
			},
		},
	}
	for _, tt := range tests {
		log.Printf("[TEST]: %v has started\n", tt.name)
		t.Run(tt.name, func(t *testing.T) {
			actual := parseFlags(tt.args.applicationCommand, tt.args.input)
			if actual != tt.expectedError {
				log.Printf("[TEST]: %v has failed\n\n", tt.name)
				t.Errorf("Failed test %v", tt.name)
			} else {
				log.Printf("[TEST]: %v has successfully finished\n\n", tt.name)
			}
		})
	}
}

func Test_checkIfFlagIsSupported(t *testing.T) {
	type args struct {
		flag string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr error
	}{
		{
			name: "flagIsNotSupported",
			args: args{
				flag: "unknown",
			},
			want:    false,
			wantErr: errUnsupportedFlag,
		},
		{
			name: "flagIsSupported",
			args: args{
				flag: "website",
			},
			want:    true,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		log.Printf("[TEST]: %v has started\n", tt.name)
		t.Run(tt.name, func(t *testing.T) {
			got, err := checkIfFlagIsSupported(tt.args.flag)
			if err != tt.wantErr {
				log.Printf("[TEST]: %v has failed\n\n", tt.name)
				t.Errorf("checkIfFlagIsSupported() error = %v, wantErr %v\n", err, tt.wantErr)
			}
			if got != tt.want {
				log.Printf("[TEST]: %v has failed\n\n", tt.name)
				t.Errorf("checkIfFlagIsSupported() = %v, want %v\n", got, tt.want)
			} else {
				log.Printf("[TEST]: %v has successfully finished\n\n", tt.name)
			}
		})
	}
}

func Test_checkIfWebsiteIsSupported(t *testing.T) {
	type args struct {
		website string
	}
	tests := []struct {
		args    args
		name    string
		want    bool
		wantErr error
	}{
		{
			args: args{
				website: "fake",
			},
			name:    "unsupportedWebsite",
			want:    false,
			wantErr: errWebsiteFlag,
		},
		{
			args: args{
				website: "amazon",
			},
			name:    "supportedWebsite",
			want:    true,
			wantErr: nil,
		},
		{
			args: args{
				website: "bestbuy",
			},
			name:    " supportedWebsite",
			want:    true,
			wantErr: nil,
		},
	}
	os.Chdir("..")
	for _, tt := range tests {
		log.Printf("[TEST]: %v has started\n", tt.name)
		t.Run(tt.name, func(t *testing.T) {
			got, err := checkIfWebsiteIsSupported(tt.args.website)
			if got != tt.want {
				log.Printf("[TEST]: %v has failed\n\n", tt.name)
				t.Errorf("checkIfWebsiteIsSupported() error = %v, wantErr %v\n", err, tt.wantErr)
			}
			if err != tt.wantErr {
				log.Printf("[TEST]: %v has failed\n\n", tt.name)
				t.Errorf("checkIfWebsiteIsSupported() = %v, want %v\n", got, tt.want)
			} else {
				log.Printf("[TEST]: %v has successfully finished\n\n", tt.name)
			}
		})
	}
}
