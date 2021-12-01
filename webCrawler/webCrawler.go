package webcrawler

import (
	"log"

	webscraper "github.com/cody6750/codywebapi/webCrawler/webScraper"
)

//WebCrawler ...
type WebCrawler struct {
	scraper        *webscraper.WebScraper
	listOfWebsites []string
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
func (w WebCrawler) Crawl(url string, depth int, itemsToget []webscraper.ScrapeItemConfiguration, ScrapeConfiguration ...webscraper.ScrapeURLConfiguration) ([]string, error) {
	list, _ := w.scraper.Scrape(url, itemsToget, ScrapeConfiguration...)
	depth--
	if depth > 0 {
		for _, url := range list {
			result, _ := w.Crawl(url, depth, itemsToget, ScrapeConfiguration...)
			w.listOfWebsites = append(w.listOfWebsites, result...)
		}
	} else {
		log.Print(w.listOfWebsites)
	}
	return list, nil
}
