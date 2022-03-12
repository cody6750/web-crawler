package handler

import (
	webcrawler "github.com/cody6750/web-crawler/pkg"
	"github.com/sirupsen/logrus"
)

var (
	Identifier string = "crawler"
)

// Crawler handler for getting items from web crawler
type Crawler struct {
	crawler    *webcrawler.WebCrawler
	logger     *logrus.Logger
	Identifier string
}

// NewCrawler returns a new crawler handler with the given logger
func NewCrawler(l *logrus.Logger) *Crawler {
	return &Crawler{Identifier: Identifier, logger: l, crawler: webcrawler.NewCrawler()}
}
