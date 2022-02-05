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

	options "github.com/cody6750/web-crawler/pkg/options"
	webcrawler "github.com/cody6750/web-crawler/pkg/webScraper"
	webscraper "github.com/cody6750/web-crawler/pkg/webScraper"
)

var (
	errorMaxDepthBelowZero = errors.New("max Depth cannot be below 0, exiting")
	errorHealthCheck       = errors.New("Health check has failed")
	errorLivenessCheck     = errors.New("Liveness check has failed")
)

// Metrics ...
type Metrics struct {
	duplicatedUrlsFound int
	urlsFound           int
	urlsVisited         int
	itemsFound          int
}

//WebCrawler ...
type WebCrawler struct {
	rootURL                   string
	pendingUrlsToCrawlCount   chan int
	pendingUrlsToCrawl        chan *webscraper.URL
	urlsToCrawl               chan *webscraper.URL
	collectWebScraperResponse chan *webscraper.ScrapeResposne
	errs                      chan error
	webScraperResponses       []*webscraper.ScrapeResposne
	stop                      chan struct{}
	metrics                   Metrics
	visited                   map[string]struct{}
	webScrapers               map[int]*webscraper.WebScraper
	wg                        sync.WaitGroup
	mapLock                   sync.Mutex
	metricsLock               sync.Mutex
	options                   *options.Options
}

//New ...
func NewCrawler() *WebCrawler {
	return NewWithOptions(options.New())
}

//NewWithOptions ...
func NewWithOptions(options *options.Options) *WebCrawler {
	crawler := &WebCrawler{}
	crawler.options = options
	return crawler
}

//Init ...
func (w *WebCrawler) init(rootURL string) error {
	w.pendingUrlsToCrawlCount = make(chan int)
	w.pendingUrlsToCrawl = make(chan *webscraper.URL)
	w.collectWebScraperResponse = make(chan *webscraper.ScrapeResposne)
	w.errs = make(chan error)
	w.urlsToCrawl = make(chan *webscraper.URL)
	w.stop = make(chan struct{}, 30)
	w.visited = make(map[string]struct{})
	w.webScrapers = make(map[int]*webscraper.WebScraper)
	w.wg = *new(sync.WaitGroup)
	w.rootURL = rootURL

	err := w.initRobotsTxtRestrictions(rootURL)
	if err != nil {
		return err
	}
	return nil

}

//Crawl ...
func (w *WebCrawler) Crawl(url string, itemsToget []webscraper.ScrapeItemConfiguration, ScrapeURLConfiguration ...webscraper.ScrapeURLConfiguration) ([]*webcrawler.ScrapeResposne, error) {
	err := w.init(url)
	wgDone := make(chan bool)
	if err != nil {
		return nil, err
	}
	if w.options.MaxDepth < 0 {
		return nil, errorMaxDepthBelowZero
	}

	go func() {
		w.processScrapedUrls([]*webscraper.URL{{RootURL: url, CurrentURL: url, CurrentDepth: 0, MaxDepth: w.options.MaxDepth}})
	}()

	go w.processCrawledLinks()

	go w.monitorCrawling()

	go w.processSrapedResponse()

	for i := 0; i < w.options.WebScraperWorkerCount; i++ {
		w.wg.Add(1)
		go func(scraperNumber int) {
			defer w.wg.Done()
			w.runWebScraper(scraperNumber, itemsToget, ScrapeURLConfiguration...)
		}(i)
	}

	err = w.readinessCheck()
	if err != nil {
		log.Print("Readiness check failed")
		return nil, err
	}

	go w.livenessCheck()

	go func() {
		w.wg.Wait()
		log.Print("Closing go routines")
		close(wgDone)
	}()
	select {
	case <-wgDone:
		// carry on
		break
	case err := <-w.errs:
		return w.webScraperResponses, err
	}
	log.Printf("Finished crawling %v", url)
	return w.webScraperResponses, nil
}

func (w *WebCrawler) runWebScraper(scraperNumber int, itemsToget []webscraper.ScrapeItemConfiguration, ScrapeURLConfiguration ...webscraper.ScrapeURLConfiguration) (*webscraper.WebScraper, error) {
	wg := new(sync.WaitGroup)
	webscraper := &webscraper.WebScraper{
		RootURL:             w.rootURL,
		ScraperNumber:       scraperNumber,
		Stop:                w.stop,
		WaitGroup:           *wg,
		BlackListedURLPaths: w.options.BlacklistedURLPaths,
		HeaderKey:           w.options.HeaderKey,
		HeaderValue:         w.options.HeaderValue,
	}

	w.mapLock.Lock()
	w.webScrapers[scraperNumber] = webscraper
	w.mapLock.Unlock()
	for {
		select {
		case url := <-w.urlsToCrawl:
			// Options: Set delay
			if w.options.CrawlDelay != 0 {
				time.Sleep(time.Second * time.Duration(w.options.CrawlDelay))
			}
			// Options: Ability to delay the execution of scraper
			if numGoRoutine := runtime.NumGoroutine(); numGoRoutine > w.options.MaxGoRoutines {
				return webscraper, fmt.Errorf("webscraper gorutines has supressed the max go routines. Current: %v Max: %v", numGoRoutine, w.options.MaxGoRoutines)
			}
			// Options: Ability to cap the number of urls scraped
			if w.options.MaxVisitedUrls <= w.metrics.urlsVisited {
				return webscraper, fmt.Errorf("url visited has supressed the max url visited. Current: %v Max: %v", w.metrics.urlsVisited, w.options.MaxVisitedUrls)
			}
			webscraper.WaitGroup.Add(1)
			go func() {
				defer webscraper.WaitGroup.Done()
				scrapeResponse, _ := webscraper.Scrape(url, itemsToget, ScrapeURLConfiguration...)
				w.incrementMetrics(&Metrics{urlsFound: len(scrapeResponse.ExtractedURLs), urlsVisited: 1, itemsFound: len(scrapeResponse.ExtractedItem)})
				log.Printf("Go routine:%v | Crawling link: %v | Current depth: %v | Url Visited: %v | Url Found : %v | Duplicate Url found: %v | Items Found: %v", scraperNumber, url.CurrentURL, url.CurrentDepth, w.metrics.urlsVisited, w.metrics.urlsFound, w.metrics.duplicatedUrlsFound, w.metrics.itemsFound)
				w.processScrapedUrls(scrapeResponse.ExtractedURLs)
				w.pendingUrlsToCrawlCount <- -1
				if !w.options.AllowEmptyItem && len(scrapeResponse.ExtractedItem) == 0 {
					return
				}
				w.collectWebScraperResponse <- scrapeResponse
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
	if metrics.duplicatedUrlsFound != 0 {
		w.metrics.duplicatedUrlsFound += metrics.duplicatedUrlsFound
	}
	w.metricsLock.Unlock()
	return metrics
}

func (w *WebCrawler) processScrapedUrls(scrapedUrls []*webscraper.URL) {
	if len(scrapedUrls) == 0 {
		return
	}
	if scrapedUrls[0].CurrentDepth > w.options.MaxDepth {
		//contiue
	} else {
		for _, url := range scrapedUrls {
			w.pendingUrlsToCrawl <- url
			w.pendingUrlsToCrawlCount <- 1
		}
	}
}

func (w *WebCrawler) processSrapedResponse() {
	for response := range w.collectWebScraperResponse {
		w.webScraperResponses = append(w.webScraperResponses, response)
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
				// log.Print("empty url")
				w.pendingUrlsToCrawlCount <- -1
				w.incrementMetrics(&Metrics{duplicatedUrlsFound: 1})
				continue
			}
			_, visited := w.visited[url.CurrentURL]
			if !visited {
				w.visited[url.CurrentURL] = struct{}{}
				w.urlsToCrawl <- url
			} else {
				// log.Print("already visited url")
				w.incrementMetrics(&Metrics{duplicatedUrlsFound: 1})
				w.pendingUrlsToCrawlCount <- -1
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
			w.stopAllWebScrapers()
			return fmt.Errorf("failed liveness check, number of go routines: %v is below threshold: %v", numGoRoutine, w.options.WebScraperWorkerCount*2)
		}
		if checkCounter {
			go func() {
				checkCounter = false
				pastCounter := w.metrics.urlsVisited
				time.Sleep(time.Second * 60)
				presentCounter := w.metrics.urlsVisited
				if pastCounter == presentCounter {
					w.stopAllWebScrapers()
					//return fmt.Errorf("Failed liveness check, url has not been crawled during 30 second interval")
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
	resp := webscraper.ConnectToWebsite(url, w.options.HeaderKey, w.options.HeaderValue)
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
