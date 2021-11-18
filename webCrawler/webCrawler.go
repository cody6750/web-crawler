package webcrawler

import (
	webscraper "github.com/cody6750/codywebapi/webCrawler/webScraper"
)

//WebCrawler ...
type WebCrawler struct {
	scraper *webscraper.WebScraper
}

//New ...
func New() *WebCrawler {

	scraper := webscraper.New()
	crawler := &WebCrawler{
		scraper: scraper,
	}

	return crawler
}

//Crawl ...
func (w WebCrawler) Crawl(url string, ScrapeConfiguration ...webscraper.ScrapeConfiguration) ([]string, error) {
	list, _ := w.scraper.Scrape(url, ScrapeConfiguration...)
	return list, nil
}

//Run ...
func Run() {
	crawl := New()
	crawl.Crawl("https://www.amazon.com/s?k=RTX+3080&ref=nb_sb_noss_2")
}
