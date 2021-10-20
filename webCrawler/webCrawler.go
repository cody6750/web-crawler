package webcrawler

import (
	"net/http"
	"time"

	"github.com/cody6750/webCrawler/webScraper"
)

//WebCrawler ...
type WebCrawler struct {
	client *http.Client
}

//New ...
func (WebCrawler) New() *WebCrawler {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	crawler := &WebCrawler{
		client: client,
	}

	return crawler
}

//Crawl ...
func (WebCrawler) Crawl(url string, scraper webScraper.WebScraper) {
}
