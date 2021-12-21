package webcrawler

import (
	"errors"
	"log"
	"sync"

	webscraper "github.com/cody6750/codywebapi/webCrawler/webScraper"
)

//WebCrawler ...
type WebCrawler struct {
	scraper        *webscraper.WebScraper
	listOfWebsites []string
	name           string
	// TODO: Turn this into []webscraper.URL
	queue              chan []string
	stop               chan struct{}
	linkVisitedCounter chan int
	VisitedLinkCounter int
	visited            map[string]struct{}
	webScrapers        map[int]*webscraper.WebScraper
	wg                 sync.WaitGroup
}

//New ...
// TODO: ADD OPTIONS. New will create a webcrawler with default options
func New() *WebCrawler {
	crawler := &WebCrawler{}
	return crawler
}

//NewWithOptions ...
// TODO: ADDO OPTIONS. New with options will create a webcrawler with overrided options.
func NewWithOptions() *WebCrawler {
	crawler := &WebCrawler{}
	return crawler
}

//Init ...
func (w *WebCrawler) Init() {
	w.scraper = webscraper.New()
	w.wg = *new(sync.WaitGroup)
	//w.queue = make(chan []webscraper.URL)
	w.queue = make(chan []string, 1)
	w.stop = make(chan struct{}, 5)
	w.visited = make(map[string]struct{})
	w.webScrapers = make(map[int]*webscraper.WebScraper)
}

//Crawl ...
func (w *WebCrawler) Crawl(url string, maxDepth int, itemsToget []webscraper.ScrapeItemConfiguration, ScrapeConfiguration ...webscraper.ScrapeURLConfiguration) ([]string, error) {
	w.Init()
	w.enqueue(url)
	if maxDepth < 0 {
		return w.listOfWebsites, errors.New("max Depth cannot be below 0, exiting")
	}
	if maxDepth == 0 {
		return append(w.listOfWebsites, url), nil
	}
	w.listOfWebsites, _ = w.scraper.Scrape(url, itemsToget, ScrapeConfiguration...)
	w.enqueueList(w.listOfWebsites)
	scraperCount := 20
	if maxDepth <= 1 {
		for i := 0; i < scraperCount; i++ {
			w.wg.Add(1)
			go func(scraperNumber int, itemsToget []webscraper.ScrapeItemConfiguration, ScrapeConfiguration ...webscraper.ScrapeURLConfiguration) {
				defer w.wg.Done()
				scraper, _ := w.launchWebScraper(scraperNumber, itemsToget, ScrapeConfiguration...)
				w.webScrapers[scraperNumber] = scraper
			}(i, itemsToget, ScrapeConfiguration...)
		}
		w.wg.Wait()
	}
	return w.listOfWebsites, nil
}

func (w *WebCrawler) launchWebScraper(scraperNumber int, itemsToget []webscraper.ScrapeItemConfiguration, ScrapeConfiguration ...webscraper.ScrapeURLConfiguration) (*webscraper.WebScraper, error) {
	wgg := new(sync.WaitGroup)
	webscraper := &webscraper.WebScraper{
		Host:          w.name,
		ScraperNumber: scraperNumber,
		Queue:         w.queue,
		Stop:          w.stop,
		Visited:       w.visited,
		Wg:            *wgg,
	}
	log.Print(scraperNumber)
	for {
		select {
		case urlsToParse := <-w.queue:
			for _, urls := range urlsToParse {
				// if w.VisitedLinkCounter > 20 {
				// 	break
				// }
				if _, visited := w.visited[urls]; !visited {
					webscraper.Wg.Add(1)
					go func(urls string, scrapNumber int) {
						defer webscraper.Wg.Done()
						w.VisitedLinkCounter++
						log.Printf("Go routine:%v | Crawling link: %v | Counter Link: %v", scraperNumber, urls, w.VisitedLinkCounter)
						scraperWebsiteList, _ := w.scraper.Scrape(urls, itemsToget, ScrapeConfiguration...)
						w.enqueueList(scraperWebsiteList)
						w.processURL(urls)
					}(urls, scraperNumber)
				}
			}
			webscraper.Wg.Wait()
			webscraper.Stop <- struct{}{}
		case <-webscraper.Stop:
			log.Print("exiting webscraper")
			return webscraper, nil
		}
	}
}

func (w *WebCrawler) processURL(url string) {
	w.visited[url] = struct{}{}
}

/*
Unbuffered Channel:
In an unbuffered channel, the send and recieve must be ready at the same time or the channel is blocked. The select statement can only execute each case statement once thus the send is never able to recieve.

Bufferd Channel:
In an buffered channel, the send and recieve are able to still send without the other operation being ready due to the buffer/capacity of the channel. It is able to hold that operation until it is called. Thus The select statement can execute.
https://stackoverflow.com/questions/47525250/in-the-go-select-construct-can-i-have-send-and-receive-to-unbuffered-channel-in
*/

func (w *WebCrawler) enqueueList(list []string) {
	toStack := list
	for {
		select {
		case w.queue <- toStack:
			return
		case oldStack := <-w.queue:
			toStack = append(oldStack, toStack...)
		}
	}
}

func (w *WebCrawler) enqueue(url string) {
	var toStack []string
	toStack = append(toStack, url)
	w.queue <- toStack
}
