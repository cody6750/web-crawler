package webcrawler

import (
	"net/http"
	"time"
)

//ConnectToWebsite Executes a HTTP request to the url and returns the response.
func ConnectToWebsite(url, headerKey, headerValue string) (*http.Response, error) {
	client := &http.Client{
		Timeout: 60 * time.Second,
	}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return request.Response, err
	}
	request.Header.Set(headerKey, headerValue)

	response, err := client.Do(request)
	if err != nil {
		return request.Response, err
	}
	return response, nil
}
