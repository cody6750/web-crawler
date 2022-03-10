package webcrawler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	options "github.com/cody6750/web-crawler/pkg/options"
	services "github.com/cody6750/web-crawler/pkg/services/aws"
	webscraper "github.com/cody6750/web-crawler/pkg/webScraper"
	"github.com/sirupsen/logrus"
)

// Metrics ...
type Metrics struct {
	URL                 string
	DuplicatedUrlsFound int
	UrlsFound           int
	UrlsVisited         int
	ItemsFound          int
}

//WebCrawler ...
type WebCrawler struct {
	rootURL                   string
	pendingUrlsToCrawlCount   chan int
	pendingUrlsToCrawl        chan *webscraper.URL
	urlsToCrawl               chan *webscraper.URL
	collectWebScraperResponse chan *webscraper.Response
	errs                      chan error
	webScraperResponses       []*webscraper.Response
	stop                      chan struct{}
	metrics                   Metrics
	visited                   map[string]struct{}
	webScrapers               map[int]*webscraper.WebScraper
	wg                        sync.WaitGroup
	mapLock                   sync.Mutex
	metricsLock               sync.Mutex
	Logger                    *logrus.Logger
	Options                   *options.Options
	// session established a session with AWS. Requires AWS to be configured on the
	// machine. The session is created through initAWS which is set using options.AWSMaxRetries
	// and options.AWSRegion or AWS_MAX_RETRIES and AWS_REGION environent variables.
	session *session.Session

	// s3Svc establishes a session with AWS S3 manager using the AWS session.
	// Allows us to upload files to S3.
	s3Svc *s3manager.Uploader
}

//Response ...
type Response struct {
	WebScraperResponses []*webscraper.Response
	Metrics             *Metrics
}

//NewCrawler ...
func NewCrawler() *WebCrawler {
	return NewWithOptions(options.New())
}

//NewWithOptions ...
func NewWithOptions(options *options.Options) *WebCrawler {
	wc := &WebCrawler{}
	wc.Logger = logrus.New()
	wc.Logger.SetFormatter(&logrus.TextFormatter{ForceColors: true, FullTimestamp: true})
	wc.Options = options
	wc.getEnvVariables()
	return wc
}

//Init ...
func (wc *WebCrawler) init(rootURL string) error {
	wc.pendingUrlsToCrawlCount = make(chan int)
	wc.pendingUrlsToCrawl = make(chan *webscraper.URL)
	wc.collectWebScraperResponse = make(chan *webscraper.Response)
	wc.errs = make(chan error)
	wc.urlsToCrawl = make(chan *webscraper.URL)
	wc.stop = make(chan struct{}, 30)
	wc.visited = make(map[string]struct{})
	wc.webScrapers = make(map[int]*webscraper.WebScraper)
	wc.wg = *new(sync.WaitGroup)
	wc.rootURL = rootURL
	err := wc.initRobotsTxtRestrictions(rootURL)
	if err != nil {
		wc.Logger.WithField("URL: ", rootURL).Info("robots.txt does not exist for website")
	}
	return nil

}

// initAWS creates the required AWS session and services.
func (wc *WebCrawler) initAWS(maxRetries int, region string) {
	configs := aws.Config{
		Region:     aws.String(region),
		MaxRetries: aws.Int(maxRetries),
	}
	wc.session = session.Must(session.NewSession(&configs))
	wc.s3Svc = s3manager.NewUploader(wc.session)
}

//Crawl ...
func (wc *WebCrawler) Crawl(url string, itemsToget []webscraper.ScrapeItemConfig, urlsToGet ...webscraper.ScrapeURLConfig) (*Response, error) {
	wc.Logger.WithField("url", url).Info("Starting to crawl url")
	wgDone := make(chan bool)
	err := wc.init(url)

	if err != nil {
		wc.Logger.WithError(err).Error("cannot initialize crawler")
		return nil, err
	}

	if wc.Options.MaxDepth < 0 {
		return nil, fmt.Errorf("max depth is cannot be lower then 0. Current max depth: %v", wc.Options.MaxDepth)
	}

	go func() {
		wc.processScrapedUrls([]*webscraper.URL{{RootURL: url, CurrentURL: url, CurrentDepth: 0, MaxDepth: wc.Options.MaxDepth}})
	}()

	go wc.processCrawledLinks()

	go wc.monitorCrawling()

	go wc.processSrapedResponse()

	wc.Logger.WithField("Webscraper Count: ", wc.Options.WebScraperWorkerCount).Debug("Deploying webscrapers")
	for i := 0; i < wc.Options.WebScraperWorkerCount; i++ {
		wc.wg.Add(1)
		go func(scraperNumber int) {
			defer wc.wg.Done()
			wc.runWebScraper(scraperNumber, itemsToget, urlsToGet...)
		}(i)
	}
	wc.Logger.WithField("Webscraper Count: ", wc.Options.WebScraperWorkerCount).Debug("Successfully deployed webscrapers")

	err = wc.readinessCheck()
	if err != nil {
		wc.Logger.WithError(err).Error("Readiness check failed")
		return nil, err
	}

	go wc.livenessCheck()

	go func() {
		wc.Logger.Debug("Waiting for all scraper goroutines to finish scraping")
		wc.wg.Wait()
		close(wgDone)
	}()
	select {
	case <-wgDone:
		// carry on
		break
	case err := <-wc.errs:
		return &Response{WebScraperResponses: wc.webScraperResponses, Metrics: &wc.metrics}, err
	}

	response := &Response{WebScraperResponses: wc.webScraperResponses, Metrics: &wc.metrics}
	if wc.Options.WriteOutputToS3 {
		wc.initAWS(wc.Options.AWSMaxRetries, wc.Options.AWSRegion)
		out, err := json.Marshal(response)
		if err != nil {
			wc.Logger.WithError(err).Error("Unable to marshal json")
			return response, err
		}
		err = services.WriteToS3(wc.s3Svc, strings.NewReader(string(out)), wc.Options.AWSS3Bucket, services.GenerateFileName("crawl_results", ".json"), "")
		if err != nil {
			wc.Logger.WithError(err).Error("Unable to upload file to S3")
			return response, err
		}
	}
	wc.Logger.WithField("url", url).Info("Finished crawling url")
	return response, nil
}

func (wc *WebCrawler) runWebScraper(scraperNumber int, itemsToget []webscraper.ScrapeItemConfig, urlsToGet ...webscraper.ScrapeURLConfig) (*webscraper.WebScraper, error) {
	wg := new(sync.WaitGroup)
	ws := &webscraper.WebScraper{
		Logger:              wc.Logger,
		RootURL:             wc.rootURL,
		ScraperNumber:       scraperNumber,
		Stop:                wc.stop,
		WaitGroup:           *wg,
		BlackListedURLPaths: wc.Options.BlacklistedURLPaths,
		HeaderKey:           wc.Options.HeaderKey,
		HeaderValue:         wc.Options.HeaderValue,
	}

	wc.mapLock.Lock()
	wc.webScrapers[scraperNumber] = ws
	wc.mapLock.Unlock()
	for {
		select {
		// If there is a url to crawl, begin scraping concurrently
		case url := <-wc.urlsToCrawl:
			// Options: Delay duration between each crawl
			if wc.Options.CrawlDelay != 0 {
				time.Sleep(time.Second * time.Duration(wc.Options.CrawlDelay))
			}
			// Options: Set maximum amount of GoRoutines. Each webscraper deploys a gorotuine per each url in the channel.
			if numGoRoutine := runtime.NumGoroutine(); numGoRoutine > wc.Options.MaxGoRoutines {
				wc.urlsToCrawl <- url
				return ws, fmt.Errorf("webscraper gorutines has supressed the max go routines. Current: %v Max: %v", numGoRoutine, wc.Options.MaxGoRoutines)
			}
			// Options: Ability to cap the number of urls scraped. Shared value between each webscraper.
			if wc.Options.MaxVisitedUrls <= wc.metrics.UrlsVisited {
				return ws, fmt.Errorf("url visited has supressed the max url visited. Current: %v Max: %v", wc.metrics.UrlsVisited, wc.Options.MaxVisitedUrls)
			}

			// Begin scraping concurrently
			ws.WaitGroup.Add(1)
			go func() {
				defer ws.WaitGroup.Done()
				scrapeResponse, err := ws.Scrape(url, itemsToget, urlsToGet...)
				if err != nil {
					wc.errs <- err
				}
				wc.incrementMetrics(&Metrics{URL: wc.rootURL, UrlsFound: len(scrapeResponse.ExtractedURLs), UrlsVisited: 1, ItemsFound: len(scrapeResponse.ExtractedItem)})
				wc.Logger.Infof("Go routine:%v | Crawling link: %v | Current depth: %v | Url Visited: %v | Url Found : %v | Duplicate Url found: %v | Items Found: %v", scraperNumber, url.CurrentURL, url.CurrentDepth, wc.metrics.UrlsVisited, wc.metrics.UrlsFound, wc.metrics.DuplicatedUrlsFound, wc.metrics.ItemsFound)
				wc.processScrapedUrls(scrapeResponse.ExtractedURLs)
				wc.pendingUrlsToCrawlCount <- -1
				if !wc.Options.AllowEmptyItem && len(scrapeResponse.ExtractedItem) == 0 {
					return
				}
				wc.collectWebScraperResponse <- scrapeResponse
			}()
		// Stop scraping, wait for all scrapes to finish before exiting function.
		case <-ws.Stop:
			ws.WaitGroup.Wait()
			return ws, nil
		}
	}
}

func (wc *WebCrawler) incrementMetrics(m *Metrics) *Metrics {
	wc.metricsLock.Lock()
	if m.URL != "" {
		wc.metrics.URL = m.URL
	}

	if m.UrlsFound != 0 {
		wc.metrics.UrlsFound += m.UrlsFound
	}

	if m.UrlsVisited != 0 {
		wc.metrics.UrlsVisited += m.UrlsVisited
	}

	if m.ItemsFound != 0 {
		wc.metrics.ItemsFound += m.ItemsFound
	}

	if m.DuplicatedUrlsFound != 0 {
		wc.metrics.DuplicatedUrlsFound += m.DuplicatedUrlsFound
	}
	wc.metricsLock.Unlock()
	return m
}

func (wc *WebCrawler) processScrapedUrls(scrapedUrls []*webscraper.URL) {
	if len(scrapedUrls) == 0 {
		return
	}

	if scrapedUrls[0].CurrentDepth > wc.Options.MaxDepth {
		//continue
	} else {
		for _, url := range scrapedUrls {
			wc.pendingUrlsToCrawl <- url
			wc.pendingUrlsToCrawlCount <- 1
		}
	}
}

func (wc *WebCrawler) processSrapedResponse() {
	for response := range wc.collectWebScraperResponse {
		wc.webScraperResponses = append(wc.webScraperResponses, response)
	}
}

//processCrawledLinks ...
func (wc *WebCrawler) processCrawledLinks() {
	for {
		url := <-wc.pendingUrlsToCrawl
		if url == nil {
			wc.Logger.Debug("Channel is closed, no longer processing crawled links")
			return
		}

		if url.CurrentURL == "" {
			wc.pendingUrlsToCrawlCount <- -1
			wc.incrementMetrics(&Metrics{DuplicatedUrlsFound: 1})
			continue
		}

		_, visited := wc.visited[url.CurrentURL]
		if !visited {
			wc.visited[url.CurrentURL] = struct{}{}
			wc.urlsToCrawl <- url
		} else {
			wc.incrementMetrics(&Metrics{DuplicatedUrlsFound: 1})
			wc.pendingUrlsToCrawlCount <- -1
		}
	}
}

//monitorCrawling ...
func (wc *WebCrawler) monitorCrawling() {
	var c int
	for count := range wc.pendingUrlsToCrawlCount {
		c += count
		if c == 0 {
			wc.Logger.Debug("No more pending urls to crawl, close all channels")
			wc.stopAllWebScrapers()
			close(wc.urlsToCrawl)
			close(wc.pendingUrlsToCrawl)
			close(wc.pendingUrlsToCrawlCount)
			wc.Logger.Debug("Sucessfully closed all channels")
		}
	}
}

func (wc *WebCrawler) shutDownWebScraper(s *webscraper.WebScraper) {
	defer wc.wg.Done()
	s.Stop <- struct{}{}
}

// StopWebScraper ...
func (wc *WebCrawler) StopWebScraper(scraperNumbers []int) {
	for _, scraperNumber := range scraperNumbers {
		if scraper, scraperExist := wc.webScrapers[scraperNumber]; scraperExist {
			wc.Logger.WithField("Number", scraperNumber).Debug("Stop webscraper")
			wc.wg.Add(1)
			go wc.shutDownWebScraper(scraper)
		}
	}
	wc.wg.Wait()
}

func (wc *WebCrawler) stopAllWebScrapers() {
	wc.Logger.Debug("Stoping all webscrapers")
	for _, scraper := range wc.webScrapers {
		wc.wg.Add(1)
		go wc.shutDownWebScraper(scraper)
	}
	wc.wg.Wait()
	wc.Logger.Debug("Successfully stopped all webscrapers")
}

// readinessCheck ensures that the specified number of webscraper workers have start up correctly which indicates that the crawler has started up correctly and are ready to scrape.
func (wc *WebCrawler) readinessCheck() error {
	time.Sleep(time.Second * 10)
	if len(wc.webScrapers) != wc.Options.WebScraperWorkerCount {
		wc.Logger.WithFields(logrus.Fields{"Current # of webscrapers": len(wc.webScrapers), "Required # of webscrapers": wc.Options.WebScraperWorkerCount}).Error("Failed health check, required # of webscrapers not reached")
		return fmt.Errorf("health check has failed")
	}

	wc.Logger.Info("Readiness check passed")
	return nil
}

func (wc *WebCrawler) livenessCheck() error {
	checkCounter := true
	for {
		time.Sleep(time.Second * 10)
		numGoRoutine := runtime.NumGoroutine()
		if numGoRoutine < wc.Options.WebScraperWorkerCount*2 {
			wc.stopAllWebScrapers()
			return fmt.Errorf("failed liveness check, number of go routines: %v is below threshold: %v", numGoRoutine, wc.Options.WebScraperWorkerCount*2)
		}

		if checkCounter {
			go func() {
				checkCounter = false
				pastCounter := wc.metrics.UrlsVisited
				time.Sleep(time.Second * 60)
				presentCounter := wc.metrics.UrlsVisited
				if pastCounter == presentCounter {
					wc.stopAllWebScrapers()
					// return fmt.Errorf("Failed liveness check, url has not been crawled during 30 second interval")
				}

				wc.Logger.Info("Liveness check passed")
				checkCounter = true
			}()
		}
		var stats runtime.MemStats
		runtime.ReadMemStats(&stats)
		wc.Logger.Debugf("HeapAlloc=%02fMB; Sys=%02fMB\n", float64(stats.HeapAlloc)/1024.0/1024.0, float64(stats.Sys)/1024.0/1024.0)
	}
}

func (wc *WebCrawler) initRobotsTxtRestrictions(url string) error {
	// Given URL, generate robots.txt url, and get the response.
	wc.Logger.WithField("url", url).Debugf("Initializing robots.txt restrictions")
	url = generateRobotsTxtURLPath(url)
	resp, err := webscraper.ConnectToWebsite(url, wc.Options.HeaderKey, wc.Options.HeaderValue)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Parses /robots.txt for blacklisted url path. Generates map used for checking.
	wc.Logger.WithField("url", url).Debugf("Parsing url for robots.txt restrictions")
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
				wc.Options.BlacklistedURLPaths[strings.ReplaceAll(strings.Split(rules, " ")[1], "*", "")] = struct{}{}
			}
		}
	}
	wc.Logger.WithField("url", url).Debugf("Successfully parsed url for robots.txt restrictions")
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
