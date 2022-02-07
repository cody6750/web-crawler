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
	"github.com/sirupsen/logrus"
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
	Logger                    *logrus.Logger
	Options                   *options.Options
}

//New ...
func NewCrawler() *WebCrawler {
	return NewWithOptions(options.New())
}

//NewWithOptions ...
func NewWithOptions(options *options.Options) *WebCrawler {
	crawler := &WebCrawler{}
	crawler.Logger = logrus.New()
	crawler.Logger.SetFormatter(&logrus.TextFormatter{ForceColors: true, FullTimestamp: true})
	crawler.Options = options
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
func (w *WebCrawler) Crawl(url string, itemsToget []webscraper.ScrapeItemConfiguration, urlsToGet ...webscraper.ScrapeURLConfiguration) ([]*webcrawler.ScrapeResposne, error) {
	w.Logger.WithField("url", url).Info("Starting to crawl url")
	wgDone := make(chan bool)
	err := w.init(url)
	if err != nil {
		w.Logger.WithError(err).Errorf("Cannot initialize crawler")
		return nil, err
	}

	if w.Options.MaxDepth < 0 {
		w.Logger.WithError(errorMaxDepthBelowZero).Errorf("Max depth is cannot be lower then 0. Current max depth: %v", w.Options.MaxDepth)
		return nil, errorMaxDepthBelowZero
	}

	go func() {
		w.processScrapedUrls([]*webscraper.URL{{RootURL: url, CurrentURL: url, CurrentDepth: 0, MaxDepth: w.Options.MaxDepth}})
	}()

	go w.processCrawledLinks()

	go w.monitorCrawling()

	go w.processSrapedResponse()

	w.Logger.WithField("Webscraper Count: ", w.Options.WebScraperWorkerCount).Debug("Deploying webscrapers")
	for i := 0; i < w.Options.WebScraperWorkerCount; i++ {
		w.wg.Add(1)
		go func(scraperNumber int) {
			defer w.wg.Done()
			w.runWebScraper(scraperNumber, itemsToget, urlsToGet...)
		}(i)
	}
	w.Logger.WithField("Webscraper Count: ", w.Options.WebScraperWorkerCount).Debug("Successfully deployed webscrapers")

	err = w.readinessCheck()
	if err != nil {
		w.Logger.WithError(err).Error("Readiness check failed")
		return nil, err
	}

	go w.livenessCheck()

	go func() {
		w.Logger.Debug("Waiting for all scraper goroutines to finish scraping")
		w.wg.Wait()
		close(wgDone)
	}()
	select {
	case <-wgDone:
		// carry on
		break
	case err := <-w.errs:
		w.Logger.WithError(err).Errorf("Failed to crawl url")
		return w.webScraperResponses, err
	}
	w.Logger.WithField("url", url).Info("Finished crawling url")
	return w.webScraperResponses, nil
}

func (w *WebCrawler) runWebScraper(scraperNumber int, itemsToget []webscraper.ScrapeItemConfiguration, ScrapeURLConfiguration ...webscraper.ScrapeURLConfiguration) (*webscraper.WebScraper, error) {
	wg := new(sync.WaitGroup)
	webscraper := &webscraper.WebScraper{
		RootURL:             w.rootURL,
		ScraperNumber:       scraperNumber,
		Stop:                w.stop,
		WaitGroup:           *wg,
		BlackListedURLPaths: w.Options.BlacklistedURLPaths,
		HeaderKey:           w.Options.HeaderKey,
		HeaderValue:         w.Options.HeaderValue,
	}

	w.mapLock.Lock()
	w.webScrapers[scraperNumber] = webscraper
	w.mapLock.Unlock()
	for {
		select {
		// If there is a url to crawl, begin scraping concurrently
		case url := <-w.urlsToCrawl:
			// Options: Delay duration between each crawl
			if w.Options.CrawlDelay != 0 {
				time.Sleep(time.Second * time.Duration(w.Options.CrawlDelay))
			}
			// Options: Set maximum amount of GoRoutines. Each webscraper deploys a gorotuine per each url in the channel.
			if numGoRoutine := runtime.NumGoroutine(); numGoRoutine > w.Options.MaxGoRoutines {
				return webscraper, fmt.Errorf("webscraper gorutines has supressed the max go routines. Current: %v Max: %v", numGoRoutine, w.Options.MaxGoRoutines)
			}
			// Options: Ability to cap the number of urls scraped. Shared value between each webscraper.
			if w.Options.MaxVisitedUrls <= w.metrics.urlsVisited {
				return webscraper, fmt.Errorf("url visited has supressed the max url visited. Current: %v Max: %v", w.metrics.urlsVisited, w.Options.MaxVisitedUrls)
			}

			// Begin scraping concurrently
			webscraper.WaitGroup.Add(1)
			go func() {
				defer webscraper.WaitGroup.Done()
				scrapeResponse, _ := webscraper.Scrape(url, itemsToget, ScrapeURLConfiguration...)
				w.incrementMetrics(&Metrics{urlsFound: len(scrapeResponse.ExtractedURLs), urlsVisited: 1, itemsFound: len(scrapeResponse.ExtractedItem)})
				w.Logger.Infof("Go routine:%v | Crawling link: %v | Current depth: %v | Url Visited: %v | Url Found : %v | Duplicate Url found: %v | Items Found: %v", scraperNumber, url.CurrentURL, url.CurrentDepth, w.metrics.urlsVisited, w.metrics.urlsFound, w.metrics.duplicatedUrlsFound, w.metrics.itemsFound)
				w.processScrapedUrls(scrapeResponse.ExtractedURLs)
				w.pendingUrlsToCrawlCount <- -1
				if !w.Options.AllowEmptyItem && len(scrapeResponse.ExtractedItem) == 0 {
					return
				}
				w.collectWebScraperResponse <- scrapeResponse
			}()
		// Stop scraping, wait for all scrapes to finish before exiting function.
		case <-webscraper.Stop:
			webscraper.WaitGroup.Wait()
			return webscraper, nil
		}
	}
}

func (w *WebCrawler) incrementMetrics(m *Metrics) *Metrics {
	w.metricsLock.Lock()
	if m.urlsFound != 0 {
		w.metrics.urlsFound += m.urlsFound
	}

	if m.urlsVisited != 0 {
		w.metrics.urlsVisited += m.urlsVisited
	}
	if m.itemsFound != 0 {
		w.metrics.itemsFound += m.itemsFound
	}
	if m.duplicatedUrlsFound != 0 {
		w.metrics.duplicatedUrlsFound += m.duplicatedUrlsFound
	}
	w.metricsLock.Unlock()
	return m
}

func (w *WebCrawler) processScrapedUrls(scrapedUrls []*webscraper.URL) {
	if len(scrapedUrls) == 0 {
		return
	}
	if scrapedUrls[0].CurrentDepth > w.Options.MaxDepth {
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
				w.Logger.Debug("Channel is closed, no longer processing crawled links")
				return
			}
			if url.CurrentURL == "" {
				w.pendingUrlsToCrawlCount <- -1
				w.incrementMetrics(&Metrics{duplicatedUrlsFound: 1})
				continue
			}
			_, visited := w.visited[url.CurrentURL]
			if !visited {
				w.visited[url.CurrentURL] = struct{}{}
				w.urlsToCrawl <- url
			} else {
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
			w.Logger.Debug("No more pending urls to crawl, close all channels")
			w.stopAllWebScrapers()
			close(w.urlsToCrawl)
			close(w.pendingUrlsToCrawl)
			close(w.pendingUrlsToCrawlCount)
			w.Logger.Debug("Sucessfully closed all channels")
		}
	}
}

func (w *WebCrawler) shutDownWebScraper(s *webscraper.WebScraper) {
	w.wg.Add(1)
	defer w.wg.Done()
	s.Stop <- struct{}{}
}

func (w *WebCrawler) stopWebScraper(scraperNumbers []int) {
	for _, scraperNumber := range scraperNumbers {
		if scraper, scraperExist := w.webScrapers[scraperNumber]; scraperExist {
			w.Logger.WithField("Number", scraperNumber).Debug("Stop webscraper")
			go w.shutDownWebScraper(scraper)
		}
	}
	w.wg.Wait()
}

func (w *WebCrawler) stopAllWebScrapers() {
	w.Logger.Debug("Stoping all webscrapers")
	for _, scraper := range w.webScrapers {
		go w.shutDownWebScraper(scraper)
	}
	w.wg.Wait()
	w.Logger.Debug("Successfully stopped all webscrapers")
}

// readinessCheck ensures that the specified number of webscraper workers have start up correctly which indicates that the crawler has started up correctly and are ready to scrape.
func (w *WebCrawler) readinessCheck() error {
	time.Sleep(time.Second * 10)
	if len(w.webScrapers) != w.Options.WebScraperWorkerCount {
		w.Logger.WithFields(logrus.Fields{"Current # of webscrapers": len(w.webScrapers), "Required # of webscrapers": w.Options.WebScraperWorkerCount}).Error("Failed health check, required # of webscrapers not reached")
		return errorHealthCheck
	}
	return nil
}

func (w *WebCrawler) livenessCheck() error {
	var checkCounter bool = true
	for {
		time.Sleep(time.Second * 10)
		numGoRoutine := runtime.NumGoroutine()
		if numGoRoutine < w.Options.WebScraperWorkerCount*2 {
			w.stopAllWebScrapers()
			return fmt.Errorf("failed liveness check, number of go routines: %v is below threshold: %v", numGoRoutine, w.Options.WebScraperWorkerCount*2)
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
	// Given URL, generate robots.txt url, and get the response.
	w.Logger.WithField("url", url).Debugf("Initializing robots.txt restrictions")
	url = generateRobotsTxtURLPath(url)
	resp := webscraper.ConnectToWebsite(url, w.Options.HeaderKey, w.Options.HeaderValue)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Parses /robots.txt for blacklisted url path. Generates map used for checking.
	w.Logger.WithField("url", url).Debugf("Parsing url for robots.txt restrictions")
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
				w.Options.BlacklistedURLPaths[strings.ReplaceAll(strings.Split(rules, " ")[1], "*", "")] = struct{}{}
			}
		}
	}
	w.Logger.WithField("url", url).Debugf("Successfully parsed url for robots.txt restrictions")
	return nil
}

func generateRobotsTxtURLPath(url string) string {
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		url = strings.Replace(url, strings.SplitAfterN(url, "/", 4)[3], "robots.txt", 1)
	} else {
		url = strings.Replace(url, strings.SplitAfterN(url, "/", 2)[1], "robots.txt", 1)
	}
	return url
}
