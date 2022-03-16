package webcrawler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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

// Metrics represents all exposed metrics by the web crawler
type Metrics struct {
	URL                 string
	DuplicatedUrlsFound int
	UrlsFound           int
	UrlsVisited         int
	ItemsFound          int
}

//Web Crawler represents all dependencies required to initialize the web crawler.
type WebCrawler struct {

	// errs is used to return errors that occur during the concurrent execution of the web scrapers
	// between the go routines.
	errs chan error

	// pendingUrlsToCrawlCount serves as global variable between channels, used to monitor the length of
	// the pendingUrlsToCrawl channel.
	pendingUrlsToCrawlCount chan int

	// stop serves as a signal reciever, used to stop the execution of the web crawler.
	stop chan struct{}

	// pendingUrlsToCrawlCount serves as a preprocessing stage for the urls that have been scraped.
	pendingUrlsToCrawl chan *webscraper.URL

	// urlsToCrawl servers as a stage that supplys the url to crawl. Ever webscraper worker will cosntantly
	// be listening and retriveing the url to crawl on this channel.
	urlsToCrawl chan *webscraper.URL

	// collectWebScraperResponsen aggregate the webs craping responses from each web scraper worker and is returned
	// in the crawler response.
	collectWebScraperResponse chan *webscraper.Response

	// webScrapers keeps track of all of the web scraper workers and their worker numbers. This is used to stop
	// specific webscraper workers.
	webScrapers map[int]*webscraper.WebScraper

	// visited used to keep track of all of the visited urls between the web scraper workers.
	visited map[string]struct{}

	// metrics represents all exposed metrics by the web crawler.
	metrics Metrics

	//Options represents the configurable options for the web crawler.
	Options *options.Options

	//mapLock used to block the initialization of the webScrapers map to prevent race conditions.
	mapLock sync.Mutex

	//mapLock used to block actions on the metrics object.
	metricsLock sync.Mutex

	//wg used to wait for channels in the web crawler.
	wg sync.WaitGroup

	//Logger used to log.
	Logger *logrus.Logger

	// session established a session with AWS. Requires AWS to be configured on the
	// machine. The session is created through initAWS which is set using options.AWSMaxRetries
	// and options.AWSRegion or AWS_MAX_RETRIES and AWS_REGION environent variables.
	session *session.Session

	// s3Svc establishes a session with AWS S3 manager using the AWS session.
	// Allows us to upload files to S3.
	s3Svc *s3manager.Uploader

	//webScraperResponses represents the aggregated responses from all web scrapers. Returned to the end user in
	// the web crawler response
	webScraperResponses []*webscraper.Response
}

//Response represents the response the web crawler returns to the end user.
type Response struct {
	WebScraperResponses []*webscraper.Response
	Metrics             *Metrics
}

//NewCrawler initializes a web crawler using the default options.
func NewCrawler() *WebCrawler {
	return NewWithOptions(options.New())
}

//NewWithOptions initializes a web crawler using custom options.
func NewWithOptions(options *options.Options) *WebCrawler {
	wc := &WebCrawler{}
	wc.Logger = logrus.New()
	wc.Logger.SetFormatter(&logrus.TextFormatter{ForceColors: true, FullTimestamp: true})
	wc.Options = options
	wc.getEnvVariables()
	if wc.Options.AWSWriteOutputToS3 {
		wc.initAWS(wc.Options.AWSMaxRetries, wc.Options.AWSRegion)
	}
	return wc
}

// init intializes all required channels and objects for the web crawler. It sets all of the
// robots.txt restrictions as well.
func (wc *WebCrawler) init(url string) error {
	wc.pendingUrlsToCrawlCount = make(chan int)
	wc.pendingUrlsToCrawl = make(chan *webscraper.URL)
	wc.collectWebScraperResponse = make(chan *webscraper.Response)
	wc.errs = make(chan error)
	wc.urlsToCrawl = make(chan *webscraper.URL)
	wc.stop = make(chan struct{}, 30)
	wc.visited = make(map[string]struct{})
	wc.webScrapers = make(map[int]*webscraper.WebScraper)
	wc.wg = *new(sync.WaitGroup)
	err := wc.initRobotsTxtRestrictions(url)
	if err != nil {
		wc.Logger.WithField("URL: ", url).Info("robots.txt does not exist for website")
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

// Crawl servers as the main function for the web crawler. It sets up all necessary channels needed to crawl, and initializes
// all of the web scraper workers. Once initialzied, the root url is processed and begins to feed the channels the next set of
// urls to crawl. This cycle continues until there are no more urls to crawl or there a stop signal is executed. All of the
// results are aggregated from all of the web scraper workers into a single web crawler response which is returned to the
// end user.
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

	//send initial URL
	go func() {
		wc.processScrapedUrls([]*webscraper.URL{{RootURL: url, CurrentURL: url, CurrentDepth: 0, MaxDepth: wc.Options.MaxDepth}})
	}()

	go wc.processCrawledUrls()

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
	log.Print(response)
	if wc.Options.AWSWriteOutputToS3 {
		out, err := json.Marshal(response)
		if err != nil {
			wc.Logger.WithError(err).Error("Unable to marshal json")
			return response, err
		}
		outputFile := services.GenerateFileName("crawl_results", ".json")
		err = services.WriteToS3(wc.s3Svc, strings.NewReader(string(out)), wc.Options.AWSS3Bucket, outputFile, "")
		if err != nil {
			wc.Logger.WithError(err).Error("Unable to upload file to S3")
			return response, err
		}
		wc.Logger.WithField("output File", outputFile).Info("Successfully uploaded file to S3")

	}
	wc.Logger.WithField("url", url).Info("Finished crawling url")
	return response, nil
}

// runWebScraper creates an instance of the web scraper, this represents a single web scraper worker. The web scraper
// worker actively listens to the urlsToCrawl channels for urls and begins to scrape them for urls and items. This function
// implements a variety of features which include delaying the crawl between urls, can restrict number of go routines runnning, ability
// to stop scraping onces if the max number of urls are visited, and collects metrics.
func (wc *WebCrawler) runWebScraper(scraperNumber int, itemsToget []webscraper.ScrapeItemConfig, urlsToGet ...webscraper.ScrapeURLConfig) (*webscraper.WebScraper, error) {
	wg := new(sync.WaitGroup)
	ws := &webscraper.WebScraper{
		Logger:              wc.Logger,
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
				wc.incrementMetrics(&Metrics{URL: url.RootURL, UrlsFound: len(scrapeResponse.ExtractedURLs), UrlsVisited: 1, ItemsFound: len(scrapeResponse.ExtractedItem)})
				wc.Logger.Infof("Go routine:%v | Crawling url: %v | Current depth: %v | Url Visited: %v | Url Found : %v | Duplicate Url found: %v | Items Found: %v", scraperNumber, url.CurrentURL, url.CurrentDepth, wc.metrics.UrlsVisited, wc.metrics.UrlsFound, wc.metrics.DuplicatedUrlsFound, wc.metrics.ItemsFound)
				log.Printf("Sending %v", &scrapeResponse.ExtractedItem)
				wc.processScrapedUrls(scrapeResponse.ExtractedURLs)
				log.Print("finished processing")
				wc.pendingUrlsToCrawlCount <- -1
				// if !wc.Options.AllowEmptyItem && len(scrapeResponse.ExtractedItem) == 0 {
				// 	return
				// }
				log.Printf("afrer %v", scrapeResponse.ExtractedItem)
				wc.collectWebScraperResponse <- scrapeResponse
			}()

			// Stop scraping, wait for all scrapes to finish before exiting function.
		case <-ws.Stop:
			ws.WaitGroup.Wait()
			return ws, nil
		}
	}
}

// incrementMetrics used to aggregate results from each web scraper worker and appends them to the existing metrics.
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

// processScrapedUrls checks the current depth of the url and decides whether or not
// to send the urls to the pendingUrlsToCrawl channel.
func (wc *WebCrawler) processScrapedUrls(scrapedUrls []*webscraper.URL) {
	log.Print(len(scrapedUrls))
	if len(scrapedUrls) == 0 {
		return
	}
	log.Print("passed length test")

	if scrapedUrls[0].CurrentDepth <= wc.Options.MaxDepth {
		log.Print("good depth")
		for _, url := range scrapedUrls {
			log.Print("adding to pending urls")
			wc.pendingUrlsToCrawl <- url
			wc.pendingUrlsToCrawlCount <- 1
		}
	}
	log.Print("done")
}

// processSrapedResponse aggregates web scraper responses from all web scraper workers
func (wc *WebCrawler) processSrapedResponse() {

	for {
		select {
		case response := <-wc.collectWebScraperResponse:
			log.Printf("Recieveing %v", response)
			wc.webScraperResponses = append(wc.webScraperResponses, response)
		}
	}
	for response := range wc.collectWebScraperResponse {
		log.Printf("Recieveing %v", response)
		wc.webScraperResponses = append(wc.webScraperResponses, response)
	}
}

// processCrawledUrls checks the url to see if its empty, if it's been visited, and generates metrics.
// If the url is ready to be crawled, it is then sent to the pendingUrlsToCrawl channel where the web scraper
// workers are actively listening to.
func (wc *WebCrawler) processCrawledUrls() {
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

// monitorCrawling used as a groutine that actively checks the pendingUrlsToCrawlCount channel to determine the state
// of the web scrapers. If the web scrapers have finished or halted or are stuck, then this function will gracefully stop
// all web scrapers and close all of the channels.
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

// shutDownWebScraper stops web scraper by sending a signal over the stop channel
func (wc *WebCrawler) shutDownWebScraper(s *webscraper.WebScraper) {
	defer wc.wg.Done()
	s.Stop <- struct{}{}
}

// StopWebScraper stops web scraper using scraper number
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

// stopAllWebScrapers stops all web scrapers.
func (wc *WebCrawler) stopAllWebScrapers() {
	wc.Logger.Debug("Stoping all webscrapers")
	for _, scraper := range wc.webScrapers {
		wc.wg.Add(1)
		go wc.shutDownWebScraper(scraper)
	}
	wc.wg.Wait()
	wc.Logger.Debug("Successfully stopped all webscrapers")
}

// readinessCheck ensures that the specified number of webscraper workers have start up correctly which indicates that
// the crawler has started up correctly and are ready to scrape.
func (wc *WebCrawler) readinessCheck() error {
	time.Sleep(time.Second * 10)
	if len(wc.webScrapers) != wc.Options.WebScraperWorkerCount {
		wc.Logger.WithFields(logrus.Fields{"Current # of webscrapers": len(wc.webScrapers), "Required # of webscrapers": wc.Options.WebScraperWorkerCount}).Error("Failed health check, required # of webscrapers not reached")
		return fmt.Errorf("health check has failed")
	}

	wc.Logger.Info("Readiness check passed")
	return nil
}

// livenessCheck determines if the web crawler is able to crawl urls, checks if web crawler is stuck crawling a url. It also
// checks the cpu and memory usages.
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

// initRobotsTxtRestrictions parses the given website robots.txt web page and intializes the user agent restrictions
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

// generateRobotsTxtURLPath given any url, generate the robots txt url path
func generateRobotsTxtURLPath(url string) string {
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		url = strings.Replace(url, strings.SplitAfterN(url, "/", 4)[3], "robots.txt", 1)
	} else {
		url = strings.Replace(url, strings.SplitAfterN(url, "/", 2)[1], "robots.txt", 1)
	}
	return url
}
