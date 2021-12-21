package webcrawler

import (
	"errors"
	"log"
	"sync"
	"time"

	webscraper "github.com/cody6750/codywebapi/webCrawler/webScraper"
)

var (
	errorMaxDepthBelowZero = errors.New("max Depth cannot be below 0, exiting")
)

//WebCrawler ...
type WebCrawler struct {
	listOfurl []string
	name      string
	// TODO: Turn this into []webscraper.URL
	stop chan struct{}

	urlsToCrawl             chan string
	pendingUrlsToCrawl      chan string
	pendingUrlsToCrawlCount chan int

	visitedUrlsCounter int
	urlsFoundCounter   int
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
	w.wg = *new(sync.WaitGroup)
	w.stop = make(chan struct{})
	w.visited = make(map[string]struct{})
	w.webScrapers = make(map[int]*webscraper.WebScraper)

	w.pendingUrlsToCrawl = make(chan string)
	w.urlsToCrawl = make(chan string)
	w.pendingUrlsToCrawlCount = make(chan int)
}

//Crawl ...
func (w *WebCrawler) Crawl(url string, maxDepth int, scraperWorkerCount int, itemsToget []webscraper.ScrapeItemConfiguration, ScrapeConfiguration ...webscraper.ScrapeURLConfiguration) ([]string, error) {
	w.Init()

	go func() {
		w.urlsToCrawl <- url
	}()
	if maxDepth < 0 {
		return w.listOfurl, errorMaxDepthBelowZero
	}
	if maxDepth == 0 {
		return append(w.listOfurl, url), nil
	}

	go w.ProcessCrawledLinks()

	go w.MonitorCrawling()

	for i := 0; i < scraperWorkerCount; i++ {
		w.wg.Add(1)
		go func(scraperNumber int) {
			defer w.wg.Done()
			scraper, _ := w.launchWebScraper(scraperNumber, maxDepth, itemsToget, ScrapeConfiguration...)
			w.webScrapers[scraperNumber] = scraper
		}(i)
	}
	w.wg.Wait()

	return w.listOfurl, nil
}

func (w *WebCrawler) launchWebScraper(scraperNumber int, maxDepth int, itemsToget []webscraper.ScrapeItemConfiguration, ScrapeConfiguration ...webscraper.ScrapeURLConfiguration) (*webscraper.WebScraper, error) {
	wgg := new(sync.WaitGroup)
	webscraper := &webscraper.WebScraper{
		Host:          w.name,
		ScraperNumber: scraperNumber,
		Stop:          w.stop,
		Visited:       w.visited,
		Wg:            *wgg,
	}
	delay := time.Duration(100000000)
	for {
		select {
		case url := <-w.urlsToCrawl:
			go func() {
				scrapedUrls, _ := webscraper.Scrape(url, itemsToget, ScrapeConfiguration...)
				w.processScrapedUrls(scrapedUrls)
				w.urlsFoundCounter += len(scrapedUrls)
				w.visitedUrlsCounter++
				log.Printf("Go routine:%v | Crawling link: %v | Counter Link: %v | Url Found : %v", scraperNumber, url, w.visitedUrlsCounter, w.urlsFoundCounter)
				w.pendingUrlsToCrawlCount <- -1
			}()
			if delay != 0 {
				time.Sleep(delay)
			}
		case <-w.stop:
			return webscraper, nil
		}
	}
}

func (w *WebCrawler) processScrapedUrls(scrapedUrls []string) {
	for _, url := range scrapedUrls {
		w.pendingUrlsToCrawl <- url
		w.pendingUrlsToCrawlCount <- 1
	}
}

//ProcessCrawledLinks ...
func (w *WebCrawler) ProcessCrawledLinks() {
	for {
		select {
		case url := <-w.pendingUrlsToCrawl:
			_, visited := w.visited[url]
			if !visited {
				w.visited[url] = struct{}{}
				w.urlsToCrawl <- url
			}
		}
	}
}

//MonitorCrawling ...
func (w *WebCrawler) MonitorCrawling() {
	var c int
	for count := range w.pendingUrlsToCrawlCount {
		c += count
		if c == 0 {
			close(w.urlsToCrawl)
			close(w.pendingUrlsToCrawl)
			close(w.pendingUrlsToCrawlCount)
			w.stop <- struct{}{}
		}
	}
}

/*
Unbuffered Channel:
In an unbuffered channel, the send and recieve must be ready at the same time or the channel is blocked. The select statement can only execute each case statement once thus the send is never able to recieve.

Bufferd Channel:
In an buffered channel, the send and recieve are able to still send without the other operation being ready due to the buffer/capacity of the channel. It is able to hold that operation until it is called. Thus The select statement can execute.
https://stackoverflow.com/questions/47525250/in-the-go-select-construct-can-i-have-send-and-receive-to-unbuffered-channel-in
*/
