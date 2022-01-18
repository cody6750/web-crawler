package webcrawler

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"runtime"
	"strings"
	"sync"
	"time"

	options "github.com/cody6750/codywebapi/webCrawler/options"
	webscraper "github.com/cody6750/codywebapi/webCrawler/webScraper"
)

var (
	errorMaxDepthBelowZero = errors.New("max Depth cannot be below 0, exiting")
	errorHealthCheck       = errors.New("Health check has failed")
	errorLivenessCheck     = errors.New("Liveness check has failed")
)

// Metrics ...
type Metrics struct {
	urlsFound   int
	urlsVisited int
	itemsFound  int
}

//WebCrawler ...
type WebCrawler struct {
	rootURL string
	// TODO: Turn this into []webscraper.URL
	pendingUrlsToCrawlCount chan int
	pendingUrlsToCrawl      chan *webscraper.URL
	urlsToCrawl             chan *webscraper.URL
	stop                    chan struct{}
	metrics                 Metrics
	visited                 map[string]struct{}
	webScrapers             map[int]*webscraper.WebScraper
	wg                      sync.WaitGroup
	mapLock                 sync.Mutex
	metricsLock             sync.Mutex
	test                    map[string]interface{}
	options                 *options.Options
}

//New ...
func New() *WebCrawler {
	return NewWithOptions(options.New())
}

//NewWithOptions ...
func NewWithOptions(options *options.Options) *WebCrawler {
	crawler := &WebCrawler{}
	crawler.options = options
	return crawler
}

//Init ...
func (w *WebCrawler) init() {
	w.pendingUrlsToCrawlCount = make(chan int)
	w.pendingUrlsToCrawl = make(chan *webscraper.URL)
	w.urlsToCrawl = make(chan *webscraper.URL)
	w.stop = make(chan struct{}, 30)
	w.visited = make(map[string]struct{})
	w.webScrapers = make(map[int]*webscraper.WebScraper)
	w.wg = *new(sync.WaitGroup)
}

//Crawl ...
func (w *WebCrawler) Crawl(url string, itemsToget []webscraper.ScrapeItemConfiguration, ScrapeURLConfiguration ...webscraper.ScrapeURLConfiguration) ([]string, error) {
	w.init()

	w.rootURL = url
	w.initRobotsTxtRestrictions(url)
	go func() {
		w.processScrapedUrls([]*webscraper.URL{{RootURL: url, CurrentURL: url, CurrentDepth: 0, MaxDepth: w.options.MaxDepth}})
	}()

	if w.options.MaxDepth < 0 {
		return nil, errorMaxDepthBelowZero
	}

	go w.processCrawledLinks()

	go w.monitorCrawling()

	for i := 0; i < w.options.WebScraperWorkerCount; i++ {
		w.wg.Add(1)
		go func(scraperNumber int) {
			defer w.wg.Done()
			w.runWebScraper(scraperNumber, itemsToget, ScrapeURLConfiguration...)
		}(i)
	}

	err := w.readinessCheck()
	if err != nil {
		log.Print("Readiness check failed")
		return nil, err
	}

	go func() {
		err = w.livenessCheck()
	}()
	if err != nil {
		log.Print("Liveness check failed")
		return nil, err
	}
	w.wg.Wait()
	log.Printf("Finished crawling %v", url)
	return []string{}, nil
}

func (w *WebCrawler) runWebScraper(scraperNumber int, itemsToget []webscraper.ScrapeItemConfiguration, ScrapeURLConfiguration ...webscraper.ScrapeURLConfiguration) (*webscraper.WebScraper, error) {
	wg := new(sync.WaitGroup)
	webscraper := &webscraper.WebScraper{
		RootURL:             w.rootURL,
		ScraperNumber:       scraperNumber,
		Stop:                w.stop,
		WaitGroup:           *wg,
		BlackListedURLPaths: w.options.BlacklistedURLPaths,
	}

	w.mapLock.Lock()
	w.webScrapers[scraperNumber] = webscraper
	w.mapLock.Unlock()
	for {
		select {
		case url := <-w.urlsToCrawl:
			// Options: Ability to delay the execution of scraper
			if numGoRoutine := runtime.NumGoroutine(); numGoRoutine > w.options.MaxGoRoutines {
				log.Print("too many go routines")
				return webscraper, nil
			}
			if w.options.CrawlDelay != 0 {
				time.Sleep(time.Second * time.Duration(w.options.CrawlDelay))
			}
			// Options: Ability to cap the number of urls scraped
			if w.options.MaxVisitedUrls <= w.metrics.urlsVisited {
				log.Print("Max urls hit... exiting")
				return webscraper, nil
			}
			// TODO: If scraper is timing out, stop that scraper object from scraping. Use stopWebCrawler() to stop specific instances of webcrawler
			webscraper.WaitGroup.Add(1)
			go func() {
				defer webscraper.WaitGroup.Done()
				scrapeResponse, _ := webscraper.Scrape(url, itemsToget, ScrapeURLConfiguration...)
				w.incrementMetrics(&Metrics{urlsFound: len(scrapeResponse.ExtractedURLs), urlsVisited: 1, itemsFound: len(scrapeResponse.ExtractedItem)})
				log.Printf("Go routine:%v | Crawling link: %v | Current depth: %v | Counter Link: %v | Url Found : %v | Items Found: %v", scraperNumber, url.CurrentURL, url.CurrentDepth, w.metrics.urlsVisited, w.metrics.urlsFound, w.metrics.itemsFound)
				if len(scrapeResponse.ExtractedURLs) != 0 {
					if scrapeResponse.ExtractedURLs[0].CurrentDepth <= w.options.MaxDepth {
						w.processScrapedUrls(scrapeResponse.ExtractedURLs)
					}
				}
				w.pendingUrlsToCrawlCount <- -1
			}()
		case <-webscraper.Stop:
			webscraper.WaitGroup.Wait()
			return webscraper, nil
		}
	}
}

func (w *WebCrawler) incrementMetrics(metrics *Metrics) *Metrics {
	w.metricsLock.Lock()
	if metrics.urlsFound != 0 {
		w.metrics.urlsFound += metrics.urlsFound
	}

	if metrics.urlsVisited != 0 {
		w.metrics.urlsVisited += metrics.urlsVisited
	}
	if metrics.itemsFound != 0 {
		w.metrics.itemsFound += metrics.itemsFound
	}
	w.metricsLock.Unlock()
	return metrics
}

func (w *WebCrawler) processScrapedUrls(scrapedUrls []*webscraper.URL) {
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
			if url == nil {
				log.Print("Channel is closed, closing processCrawledLinks goroutine")
				return
			}
			if url.CurrentURL == "" {
				continue
			}
			_, visited := w.visited[url.CurrentURL]
			if !visited {
				w.visited[url.CurrentURL] = struct{}{}
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
			log.Print("Closing channels")
			w.stopAllWebScrapers()
			close(w.urlsToCrawl)
			close(w.pendingUrlsToCrawl)
			close(w.pendingUrlsToCrawlCount)
		}
	}
}

func (w *WebCrawler) shutDownWebScraper(scraper *webscraper.WebScraper) {
	w.wg.Add(1)
	defer w.wg.Done()
	scraper.Stop <- struct{}{}
}

func (w *WebCrawler) stopWebScraper(scraperNumbers []int) {
	for _, scraperNumber := range scraperNumbers {
		if scraper, scraperExist := w.webScrapers[scraperNumber]; scraperExist {
			go w.shutDownWebScraper(scraper)
		}
	}
	w.wg.Wait()
}

func (w *WebCrawler) stopAllWebScrapers() {
	for _, scraper := range w.webScrapers {
		go w.shutDownWebScraper(scraper)
	}
	w.wg.Wait()
}

// readinessCheck ensures that the specified number of webscraper workers have start up correctly which indicates that the crawler has started up correctly and are ready to scrape.
func (w *WebCrawler) readinessCheck() error {
	time.Sleep(time.Second * 10)
	if len(w.webScrapers) != w.options.WebScraperWorkerCount {
		log.Printf("Failed health check, number of web scrapers: %v is below threshold: %v", len(w.webScrapers), w.options.WebScraperWorkerCount)
		return errorHealthCheck
	}
	return nil
}

func (w *WebCrawler) livenessCheck() error {
	var checkCounter bool = true
	for {
		time.Sleep(time.Second * 10)
		numGoRoutine := runtime.NumGoroutine()
		if numGoRoutine < w.options.WebScraperWorkerCount*2 {
			log.Printf("Failed liveness check, number of go routines: %v is below threshold: %v", numGoRoutine, w.options.WebScraperWorkerCount*2)
			return errorLivenessCheck
		}
		if checkCounter {
			go func() {
				checkCounter = false
				pastCounter := w.metrics.urlsVisited
				time.Sleep(time.Second * 60)
				presentCounter := w.metrics.urlsVisited
				if pastCounter == presentCounter {
					log.Fatalf("Failed liveness check, url has not been crawled during 30 second interval")
				}
				log.Print("Check counter passed")
				checkCounter = true
			}()
		}
		var stats runtime.MemStats
		runtime.ReadMemStats(&stats)
		fmt.Printf("HeapAlloc=%02fMB; Sys=%02fMB\n", float64(stats.HeapAlloc)/1024.0/1024.0, float64(stats.Sys)/1024.0/1024.0)
	}
}

func (w *WebCrawler) initRobotsTxtRestrictions(url string) error {
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		url = strings.Replace(url, strings.SplitAfterN(url, "/", 4)[3], "robots.txt", 1)
	} else {
		url = strings.Replace(url, strings.SplitAfterN(url, "/", 2)[1], "robots.txt", 1)
	}
	resp := webscraper.ConnectToWebsite(url)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	robotTxtRules := strings.Split(string(body), "\n")
	addRule := false
	for _, rules := range robotTxtRules {
		if strings.Contains(rules, "User-agent: *") {
			addRule = true
			continue
		} else if strings.Contains(rules, "User-agent: ") && addRule {
			break
		}
		if addRule {
			if strings.Contains(rules, "Disallow: ") {
				w.options.BlacklistedURLPaths[strings.ReplaceAll(strings.Split(rules, " ")[1], "*", "")] = struct{}{}
			}
		}
	}
	return nil
}

/*
Unbuffered Channel:
In an unbuffered channel, the send and recieve must be ready at the same time or the channel is blocked. The select statement can only execute each case statement once thus the send is never able to recieve.

Bufferd Channel:
In an buffered channel, the send and recieve are able to still send without the other operation being ready due to the buffer/capacity of the channel. It is able to hold that operation until it is called. Thus The select statement can execute.
https://stackoverflow.com/questions/47525250/in-the-go-select-construct-can-i-have-send-and-receive-to-unbuffered-channel-in
*/
