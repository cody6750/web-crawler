package webcrawler

import (
	"log"
	"net/http"
	"time"
)

//ConnectToWebsite ...
func ConnectToWebsite(WebPageURL string) *http.Response {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	//GET request to domain for HTML response
	request, err := http.NewRequest("GET", WebPageURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Set header as User-Agent so the server admins don't block our IP address from HTTP requests
	request.Header.Set("User-Agent", "This bot just searches amazon for a product")

	// Make HTTP request
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	return response
}
