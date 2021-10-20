package webcrawler

import (
	"net/http"
	"time"

	"github.com/cody6750/webcrawler/webCrawler/webScraper"
)

//WebCrawler ...
type WebCrawler struct {
	client  *http.Client
	scraper webScraper.WebScraper
}

//New ...
func (WebCrawler) New() *WebCrawler {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	scraper := &webScraper.WebScraper{
		Timeout: 30 * time.Second,
	}
	crawler := &WebCrawler{
		client:  client,
		scraper: scraper,
	}

	return crawler
}

//Crawl ...
func (WebCrawler) Crawl(url string) {
}
