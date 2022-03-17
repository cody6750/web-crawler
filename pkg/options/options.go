package options

var (
	defaultAllowEmptyItem        bool   = false
	defaultAWSWriteOutputToS3    bool   = false
	defaultAWSMaxRetries         int    = 5
	defaultCrawlDelay            int    = 5
	defaultMaxDepth              int    = 1
	defaultMaxGoRoutines         int    = 10000
	defaultMaxVisitedUrls        int    = 20
	defeaultMaxItemsFound        int    = 5000
	defaultWebScraperWorkercount int    = 5
	defaultAWSRegion             string = "us-east-1"
	defaultAWSS3Bucket           string = "webcrawler-results"
	defaultHeaderKey             string = "User-Agent"
	defaultHeaderValue           string = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36"
)

// Options ...
type Options struct {
	AllowEmptyItem        bool
	AWSWriteOutputToS3    bool
	AWSMaxRetries         int
	CrawlDelay            int
	MaxDepth              int
	MaxGoRoutines         int
	MaxVisitedUrls        int
	MaxItemsFound         int
	WebScraperWorkerCount int
	BlacklistedURLPaths   map[string]struct{}
	AWSRegion             string
	AWSS3Bucket           string
	HeaderKey             string
	HeaderValue           string
}

//New ...
func New() *Options {
	return &Options{
		AllowEmptyItem:        defaultAllowEmptyItem,
		AWSWriteOutputToS3:    defaultAWSWriteOutputToS3,
		AWSMaxRetries:         defaultAWSMaxRetries,
		CrawlDelay:            defaultCrawlDelay,
		MaxDepth:              defaultMaxDepth,
		MaxGoRoutines:         defaultMaxGoRoutines,
		MaxVisitedUrls:        defaultMaxVisitedUrls,
		MaxItemsFound:         defeaultMaxItemsFound,
		WebScraperWorkerCount: defaultWebScraperWorkercount,
		BlacklistedURLPaths:   map[string]struct{}{},
		HeaderKey:             defaultHeaderKey,
		AWSRegion:             defaultAWSRegion,
		AWSS3Bucket:           defaultAWSS3Bucket,
		HeaderValue:           defaultHeaderValue,
	}
}
