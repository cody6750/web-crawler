package webcrawler

import (
	"log"
	"net/http"
	"time"
)

//ConnectToWebsite ...
func ConnectToWebsite(webPageURL, headerKey, headerValue string) *http.Response {
	client := &http.Client{
		Timeout: 60 * time.Second,
	}
	request, err := http.NewRequest("GET", webPageURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	// ctx, cancel := context.WithTimeout(request.Context(), 60*time.Second)
	// defer cancel()
	request.Header.Set(headerKey, headerValue)
	// request = request.WithContext(ctx)
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	return response
}
