package handler

import (
	"github.com/sirupsen/logrus"
)

var (
	Identifier string = "crawler"
)

// Crawler ...
type Crawler struct {
	logger     *logrus.Logger
	Identifier string
}

// NewCrawler ...
func NewCrawler(l *logrus.Logger) *Crawler {
	return &Crawler{Identifier: Identifier, logger: l}
}
