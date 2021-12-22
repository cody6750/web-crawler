package webcrawler

import (
	"errors"
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	webscraper "github.com/cody6750/codywebapi/webCrawler/webScraper"
)

var (
	errorMaxDepthBelowZero = errors.New("max Depth cannot be below 0, exiting")
	errorHealthCheck       = errors.New("Health check has failed")
	errorLivenessCheck     = errors.New("Liveness check has failed")
)

//WebCrawler ...
type WebCrawler struct {
	name string
	// TODO: Turn this into []webscraper.URL
	pendingUrlsToCrawlCount chan int
	pendingUrlsToCrawl      chan string
	urlsToCrawl             chan string
	stop                    chan struct{}
	urlsFoundCounter        int
	visitedUrlsCounter      int
	visited                 map[string]struct{}
	webScrapers             map[int]*webscraper.WebScraper
	wg                      sync.WaitGroup
	mutex                   sync.Mutex
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
func (w *WebCrawler) init() {
	w.pendingUrlsToCrawlCount = make(chan int)
	w.pendingUrlsToCrawl = make(chan string)
	w.urlsToCrawl = make(chan string)
	w.stop = make(chan struct{})
	w.visited = make(map[string]struct{})
	w.webScrapers = make(map[int]*webscraper.WebScraper)
	w.wg = *new(sync.WaitGroup)
}

//Crawl ...
func (w *WebCrawler) Crawl(url string, maxDepth int, scraperWorkerCount int, itemsToget []webscraper.ScrapeItemConfiguration, ScrapeConfiguration ...webscraper.ScrapeURLConfiguration) ([]string, error) {
	w.init()

	go func() {
		w.urlsToCrawl <- url
	}()
	if maxDepth < 0 {
		return nil, errorMaxDepthBelowZero
	}
	if maxDepth == 0 {
		return append([]string{}, url), nil
	}

	go w.processCrawledLinks()

	go w.monitorCrawling()

	for i := 0; i < scraperWorkerCount; i++ {
		w.wg.Add(1)
		go func(scraperNumber int) {
			defer w.wg.Done()
			w.runWebScraper(scraperNumber, maxDepth, itemsToget, ScrapeConfiguration...)
		}(i)
	}

	err := w.readinessCheck(scraperWorkerCount)
	if err != nil {
		log.Print("Readiness check failed")
		return nil, err
	}

	go func() {
		err = w.livenessCheck(scraperWorkerCount)
	}()
	if err != nil {
		log.Print("Liveness check failed")
		return nil, err
	}
	w.wg.Wait()

	return []string{}, nil
}

func (w *WebCrawler) runWebScraper(scraperNumber int, maxDepth int, itemsToget []webscraper.ScrapeItemConfiguration, ScrapeConfiguration ...webscraper.ScrapeURLConfiguration) (*webscraper.WebScraper, error) {
	wgg := new(sync.WaitGroup)
	webscraper := &webscraper.WebScraper{
		Host:          w.name,
		ScraperNumber: scraperNumber,
		Stop:          w.stop,
		Wg:            *wgg,
	}

	w.mutex.Lock()
	w.webScrapers[scraperNumber] = webscraper
	w.mutex.Unlock()

	maxVisitedUrls := 200000
	delay := 2
	for {
		select {
		case url := <-w.urlsToCrawl:
			// Options: Ability to delay the execution of scraper
			if numGoRoutine := runtime.NumGoroutine(); numGoRoutine > 30 {
				continue
			}
			if delay != 0 {
				time.Sleep(time.Second * time.Duration(delay))
			}
			// Options: Ability to cap the number of urls scraped
			if maxVisitedUrls != 0 && maxVisitedUrls < w.visitedUrlsCounter {
				log.Print("Max urls hit... exiting")
				w.stopAllWebCrawlers()
			}
			// TODO: If scraper is timing out, stop that scraper object from scraping. Use stopWebCrawler() to stop specific instances of webcrawler
			webscraper.Wg.Add(1)
			go func() {
				defer webscraper.Wg.Done()
				scrapedUrls, _ := webscraper.Scrape(url, itemsToget, ScrapeConfiguration...)
				w.processScrapedUrls(scrapedUrls)
				w.urlsFoundCounter += len(scrapedUrls)
				w.visitedUrlsCounter++
				log.Printf("Go routine:%v | Crawling link: %v | Counter Link: %v | Url Found : %v Scraper Url : %v", scraperNumber, url, w.visitedUrlsCounter, w.urlsFoundCounter, scrapedUrls)
			}()
		case <-webscraper.Stop:
			webscraper.Wg.Wait()
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

//processCrawledLinks ...
func (w *WebCrawler) processCrawledLinks() {
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

//monitorCrawling ...
func (w *WebCrawler) monitorCrawling() {
	var c int
	for count := range w.pendingUrlsToCrawlCount {
		c += count
		if c == 0 {
			w.stopAllWebCrawlers()
			close(w.urlsToCrawl)
			close(w.pendingUrlsToCrawl)
			close(w.pendingUrlsToCrawlCount)
		}
	}
}

func (w *WebCrawler) stopWebCrawler(scraperNumbers []int) {
	for _, scraperNumber := range scraperNumbers {
		if scraper, scraperExist := w.webScrapers[scraperNumber]; scraperExist {
			w.wg.Add(1)
			go func(scraper webscraper.WebScraper) {
				defer w.wg.Done()
				scraper.Stop <- struct{}{}
			}(*scraper)
		}
	}
	w.wg.Wait()
}

func (w *WebCrawler) stopAllWebCrawlers() {
	for _, scraper := range w.webScrapers {
		w.wg.Add(1)
		go func(scraper webscraper.WebScraper) {
			defer w.wg.Done()
			scraper.Stop <- struct{}{}
		}(*scraper)
	}
	w.wg.Wait()
}

// readinessCheck ensures that the specified number of webscraper workers have start up correctly which indicates that the crawler has started up correctly and are ready to scrape.
func (w *WebCrawler) readinessCheck(scraperWorkerCount int) error {
	time.Sleep(time.Second * 20)
	if len(w.webScrapers) != scraperWorkerCount {
		log.Printf("Failed health check, number of web scrapers: %v is below threshold: %v", len(w.webScrapers), scraperWorkerCount)
		return errorHealthCheck
	}
	return nil
}

func (w *WebCrawler) livenessCheck(threshold int) error {
	threshold *= 2
	for {
		time.Sleep(time.Second * 10)
		numGoRoutine := runtime.NumGoroutine()
		if numGoRoutine < threshold {
			log.Printf("Failed health check, number of go routines: %v is below threshold: %v", numGoRoutine, threshold)
			return errorLivenessCheck
		}
		var stats runtime.MemStats
		runtime.ReadMemStats(&stats)
		fmt.Printf("HeapAlloc=%02fMB; Sys=%02fMB\n", float64(stats.HeapAlloc)/1024.0/1024.0, float64(stats.Sys)/1024.0/1024.0)
	}
}

/*
Unbuffered Channel:
In an unbuffered channel, the send and recieve must be ready at the same time or the channel is blocked. The select statement can only execute each case statement once thus the send is never able to recieve.

Bufferd Channel:
In an buffered channel, the send and recieve are able to still send without the other operation being ready due to the buffer/capacity of the channel. It is able to hold that operation until it is called. Thus The select statement can execute.
https://stackoverflow.com/questions/47525250/in-the-go-select-construct-can-i-have-send-and-receive-to-unbuffered-channel-in
*/
