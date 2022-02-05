package webcrawler

import (
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

func TestConnectToWebsite(t *testing.T) {
	type args struct {
		WebPageURL  string
		headerKey   string
		headerValue string
	}
	tests := []struct {
		name string
		args args
		want *http.Response
	}{
		{
			args: args{
				headerKey:   "User Agent",
				headerValue: "Harmless test",
				WebPageURL:  "https://www.newegg.com/p/pl?d=rtx+3080&LeftPriceRange=1000+",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := ConnectToWebsite(tt.args.WebPageURL, tt.args.headerKey, tt.args.headerValue)
			//We Read the response body on the line below.
			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				log.Fatalln(err)
			}
			//Convert the body to type string
			sb := string(body)
			log.Print(sb)
		})
	}
}
