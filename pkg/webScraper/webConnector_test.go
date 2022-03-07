package webcrawler

import (
	"net/http"
	"reflect"
	"testing"
)

func TestConnectToWebsite(t *testing.T) {
	type args struct {
		url         string
		headerKey   string
		headerValue string
	}
	tests := []struct {
		name    string
		args    args
		want    *http.Response
		wantErr bool
	}{
		{

			args:    args{url: "https://www.bhphotovideo.com/c/search?Ntt=RTX%203080&N=0&InitialSearch=yes&sts=ma"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConnectToWebsite(tt.args.url, tt.args.headerKey, tt.args.headerValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConnectToWebsite() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConnectToWebsite() = %v, want %v", got, tt.want)
			}
		})
	}
}
