package options

var (
	defaultCrawlDelay            int    = 2
	defaultMaxDepth              int    = 1
	defaultMaxGoRoutines         int    = 10000
	defaultMaxVisitedUrls        int    = 100
	defeaultMaxItemsFound        int    = 5000
	defaultWebScraperWorkercount int    = 5
	defaultHeaderKey             string = "User-Agent"
	defaultHeaderValue           string = "Simple_Web_Crawler_Used_For_Product_Searches"
)

// Options ...
type Options struct {
	CrawlDelay            int
	MaxDepth              int
	MaxGoRoutines         int
	MaxVisitedUrls        int
	MaxItemsFound         int
	WebScraperWorkerCount int
	BlacklistedURLPaths   map[string]struct{}
	HeaderKey             string
	HeaderValue           string
}

//New ...
func New() *Options {
	return &Options{
		CrawlDelay:            defaultCrawlDelay,
		MaxDepth:              defaultMaxDepth,
		MaxGoRoutines:         defaultMaxGoRoutines,
		MaxVisitedUrls:        defaultMaxVisitedUrls,
		MaxItemsFound:         defeaultMaxItemsFound,
		WebScraperWorkerCount: defaultWebScraperWorkercount,
		BlacklistedURLPaths:   map[string]struct{}{},
		HeaderKey:             defaultHeaderKey,
		HeaderValue:           defaultHeaderValue,
	}
}
