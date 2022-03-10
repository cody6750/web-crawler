package webcrawler

import (
	"os"

	env "github.com/cody6750/web-crawler/shared"
)

// getEnvVariables allows users to override the values within the tracking bot options via environment variables.
// It retrives the environment variables set at runtime and within the Docker container.
func (wc *WebCrawler) getEnvVariables() {
	var err error
	wc.Logger.Info("Getting environment variables")

	if os.Getenv("ALLOW_EMPTY_ITEM") != "" {
		wc.Options.AllowEmptyItem, err = env.GetEnvBool("ALLOW_EMPTY_ITEM")
		if err != nil {
			wc.Logger.WithError(err).Fatal("Failed to convert ALLOW_EMPTY_ITEM from string to bool")
		}
		wc.Logger.WithField("ALLOW_EMPTY_ITEM: ", wc.Options.AllowEmptyItem).Info("Successfully got environment variable")
	}

	if os.Getenv("WRITE_OUTPUT_TO_S3") != "" {
		wc.Options.WriteOutputToS3, err = env.GetEnvBool("WRITE_OUTPUT_TO_S3")
		if err != nil {
			wc.Logger.WithError(err).Fatal("Failed to convert WRITE_OUTPUT_TO_S3 from string to bool")
		}
		wc.Logger.WithField("WRITE_OUTPUT_TO_S3: ", wc.Options.WriteOutputToS3).Info("Successfully got environment variable")
	}

	if os.Getenv("AWS_MAX_RERIES") != "" {
		wc.Options.AWSMaxRetries, err = env.GetEnvInt("AWS_MAX_RERIES")
		if err != nil {
			wc.Logger.WithError(err).Fatal("Failed to convert AWS_MAX_RERIES from string to int")
		}
		wc.Logger.WithField("AWS_MAX_RERIES: ", wc.Options.AWSMaxRetries).Info("Successfully got environment variable")
	}

	if os.Getenv("CRAWL_DELAY") != "" {
		wc.Options.CrawlDelay, err = env.GetEnvInt("CRAWL_DELAY")
		if err != nil {
			wc.Logger.WithError(err).Fatal("Failed to convert  from string to int")
		}
		wc.Logger.WithField("CRAWL_DELAY: ", wc.Options.CrawlDelay).Info("Successfully got environment variable")
	}

	if os.Getenv("MAX_DEPTH") != "" {
		wc.Options.MaxDepth, err = env.GetEnvInt("MAX_DEPTH")
		if err != nil {
			wc.Logger.WithError(err).Fatal("Failed to convert MAX_DEPTH from string to int")
		}
		wc.Logger.WithField("MAX_DEPTH: ", wc.Options.MaxDepth).Info("Successfully got environment variable")
	}

	if os.Getenv("MAX_GO_ROUTINES") != "" {
		wc.Options.MaxGoRoutines, err = env.GetEnvInt("MAX_GO_ROUTINES")
		if err != nil {
			wc.Logger.WithError(err).Fatal("Failed to convert MAX_GO_ROUTINES from string to int")
		}
		wc.Logger.WithField("MAX_GO_ROUTINES: ", wc.Options.MaxGoRoutines).Info("Successfully got environment variable")
	}

	if os.Getenv("MAX_VISITED_URLS") != "" {
		wc.Options.MaxVisitedUrls, err = env.GetEnvInt("MAX_VISITED_URLS")
		if err != nil {
			wc.Logger.WithError(err).Fatal("Failed to convert MAX_VISITED_URLS from string to int")
		}
		wc.Logger.WithField("MAX_VISITED_URLS: ", wc.Options.MaxVisitedUrls).Info("Successfully got environment variable")
	}

	if os.Getenv("MAX_ITEMS_FOUND") != "" {
		wc.Options.MaxItemsFound, err = env.GetEnvInt("MAX_ITEMS_FOUND")
		if err != nil {
			wc.Logger.WithError(err).Fatal("Failed to convert MAX_ITEMS_FOUND from string to int")
		}
		wc.Logger.WithField("MAX_ITEMS_FOUND: ", wc.Options.MaxItemsFound).Info("Successfully got environment variable")
	}

	if os.Getenv("WEB_SCRAPER_WORKER_COUNT") != "" {
		wc.Options.WebScraperWorkerCount, err = env.GetEnvInt("WEB_SCRAPER_WORKER_COUNT")
		if err != nil {
			wc.Logger.WithError(err).Fatal("Failed to convert WEB_SCRAPER_WORKER_COUNT from string to int")
		}
		wc.Logger.WithField("WEB_SCRAPER_WORKER_COUNT: ", wc.Options.WebScraperWorkerCount).Info("Successfully got environment variable")
	}

	if os.Getenv("AWS_REGION") != "" {
		wc.Options.AWSRegion = os.Getenv("AWS_REGION")
		wc.Logger.WithField(": ", wc.Options.AWSRegion).Info("Successfully got environment variable")
	}

	if os.Getenv("AWS_S3_BUCKET") != "" {
		wc.Options.AWSS3Bucket = os.Getenv("AWS_S3_BUCKET")
		wc.Logger.WithField(": ", wc.Options.AWSS3Bucket).Info("Successfully got environment variable")
	}

	if os.Getenv("HEADER_KEY") != "" {
		wc.Options.HeaderKey = os.Getenv("HEADER_KEY")
		wc.Logger.WithField(": ", wc.Options.HeaderKey).Info("Successfully got environment variable")
	}

	if os.Getenv("HEADER_VALUE") != "" {
		wc.Options.HeaderValue = os.Getenv("HEADER_VALUE")
		wc.Logger.WithField(": ", wc.Options.HeaderValue).Info("Successfully got environment variable")
	}

	wc.Logger.Info("Successfully got environment variables")

}
