package options

var (
	defaultCrawlDelay            int    = 5
	defaultMaxDepth              int    = 1
	defaultMaxGoRoutines         int    = 10000
	defaultMaxVisitedUrls        int    = 10
	defeaultMaxItemsFound        int    = 5000
	defaultWebScraperWorkercount int    = 5
	defaultHeaderKey             string = "User-Agent"
	defaultHeaderValue           string = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36"
	defaultAllowEmptyItem        bool   = false
)

// Options ...
type Options struct {
	AllowEmptyItem        bool
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
		AllowEmptyItem:        defaultAllowEmptyItem,
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
