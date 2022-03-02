package webcrawler

import (
	"log"
	"net/http"
	"time"
)

//ConnectToWebsite ...
func ConnectToWebsite(url, headerKey, headerValue string) *http.Response {
	client := &http.Client{
		Timeout: 60 * time.Second,
	}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set(headerKey, headerValue)

	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	return response
}
