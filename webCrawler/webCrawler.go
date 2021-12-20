package webcrawler

import (
	"log"
	"sync"

	webscraper "github.com/cody6750/codywebapi/webCrawler/webScraper"
)

//WebCrawler ...
type WebCrawler struct {
	scraper        *webscraper.WebScraper
	listOfWebsites []string
	mame           string
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
func (w WebCrawler) Crawl(url string, duplicateUrls chan map[string]bool, depth int, itemsToget []webscraper.ScrapeItemConfiguration, ScrapeConfiguration ...webscraper.ScrapeURLConfiguration) ([]string, error) {
	var wg sync.WaitGroup
	var currentQueue Queue
	if depth == 0 {
		return append(w.listOfWebsites, url), nil
	}
	w.listOfWebsites, _ = w.scraper.Scrape(url, itemsToget, ScrapeConfiguration...)
	currentQueue.enqueueList(w.listOfWebsites)
	crawlerCount := 50
	if depth == 1 {
		for i := 0; i < crawlerCount; i++ {
			wg.Add(1)
			website := currentQueue.dequeue()
			go func(website string) {
				defer wg.Done()
				for len(currentQueue) != 0 {
					scraperWebsiteList, _ := w.scraper.Scrape(website, itemsToget, ScrapeConfiguration...)
					currentQueue.enqueueList(scraperWebsiteList)
					log.Printf("Go routine:%v | Crawling link: %v", i, website)
				}
			}(website)
		}
	}

	// if depth == 1 {
	// 	for i := 0; i < 10; i++ {
	// 		log.Print("called")
	// 		currentQueue := <-urlsToCrawl
	// 		scrape := currentQueue.dequeue()
	// 		wg.Add(1)
	// 		go func(i int, scrape string) {
	// 			defer wg.Done()
	// 			_, err := w.scraper.Scrape(scrape, itemsToget, ScrapeConfiguration...)
	// 			log.Printf("Go routine:%v | Crawling link: %v", i, scrape)
	// 			if err != nil {
	// 			}
	// 		}(i, scrape)
	// 		urlsToCrawl <- currentQueue
	// 	}
	// 	wg.Wait()
	// 	return w.listOfWebsites, nil
	// }
	//log.Printf("Crawling link: %v", url)
	// depth--
	// if depth > 0 {
	// 	for _, url := range w.listOfWebsites {
	// 		result, _ := w.Crawl(url, depth, itemsToget, ScrapeConfiguration...)
	// 		w.listOfWebsites = append(w.listOfWebsites, result...)
	// 	}
	// } else {
	// 	log.Print(w.listOfWebsites)
	// }
	return w.listOfWebsites, nil
}
