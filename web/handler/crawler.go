package handler

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

// Crawler ...
type Crawler struct {
	logger *logrus.Logger
}

// NewCrawler ...
func NewCrawler(l *logrus.Logger) *Crawler {
	return &Crawler{l}
}

func (c *Crawler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	c.logger.WithFields(logrus.Fields{"Method": http.MethodGet, "Url": r.URL}).Info("Serving HTTP request for:")
	switch {
	case r.Method == http.MethodGet:
		err := c.getItem(rw, r)
		if err != nil {
			c.logger.WithError(err).Errorf("Failed get request for %v", r.URL)
		}
		c.logger.WithFields(logrus.Fields{"Method": http.MethodGet, "Url": r.URL}).Info("Successfully served HTTP request for:")
	default:
		c.logger.WithFields(logrus.Fields{"Method": http.MethodGet, "Url": r.URL}).Error("HTTP request not supported for:")
		rw.WriteHeader(http.StatusMethodNotAllowed)
	}
	return
}
