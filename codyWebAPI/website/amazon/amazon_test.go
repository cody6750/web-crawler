package amazon

import (
	"log"
	"testing"
)

func Test_generateSearchURL(t *testing.T) {
	type args struct {
		item  string
		items []string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr error
	}{
		{
			name: "supported search",
			args: args{
				item: "RTX 3080",
			},
			want:    "https://www.amazon.com/s?k=RTX+3080&ref=nb_sb_noss_2",
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		log.Printf("[TEST]: %v has started\n", tt.name)
		t.Run(tt.name, func(t *testing.T) {
			got, err := generateSearchURL(tt.args.item)
			if err != tt.wantErr {
				log.Printf("[TEST]: %v has failed want: %v got: %v\n\n", tt.name, err, tt.wantErr)
				t.Errorf("generateSearchURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				log.Printf("[TEST]: %v has failed want: %v got: %v\n\n", tt.name, tt.want, got)
				t.Errorf("generateSearchURL() error = %v, got: %v", tt.want, got)
				return
			}
			log.Printf("[TEST]: %v has successfully finished\n\n", tt.name)
		})
	}
}

func TestAmazon_SearchWebsite(t *testing.T) {
	type args struct {
		item string
	}
	tests := []struct {
		name         string
		amazonObject Amazon
		args         args
		want         error
	}{
		{
			name:         "search for RTX",
			amazonObject: Amazon{Name: "6"},
			args:         args{item: "RTX 3080"},
			want:         nil,
		},
		// {
		// 	name:         "search for snowboard",
		// 	amazonObject: Amazon{Name: "2"},
		// 	args:         args{item: "Girl Clothes"},
		// 	want:         nil,
		// },
	}
	for _, tt := range tests {
		log.Printf("[TEST]: %v has started\n", tt.name)
		t.Run(tt.name, func(t *testing.T) {
			_, got := tt.amazonObject.SearchWebsite(tt.args.item)
			if tt.want != got {
				log.Printf("[TEST]: %v has failed want: %v got: %v\n\n", tt.name, tt.want, got)
				t.Errorf("generateSearchURL() error = %v, got: %v", tt.want, got)
				return
			}
			log.Printf("[TEST]: %v has successfully finished\n\n", tt.name)
		})
	}
}
