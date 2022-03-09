package handler

import (
	"context"
	"net/http"

	"github.com/cody6750/web-crawler/web/data"
	"github.com/sirupsen/logrus"
)

func (c *Crawler) GetItem(rw http.ResponseWriter, r *http.Request) {
	c.logger.WithFields(logrus.Fields{"Handler": c.Identifier, "Function": "getItem"}).Info("Starting to call handler")
	payload := r.Context().Value(KeyItem{}).(data.Payload)
	products, err := data.GetItem(c.logger, payload.RootURL, payload.ScrapeItemConfiguration, payload.ScrapeURLConfiguration...)
	if err != nil {
		c.logger.WithError(err).Error("Unable to call GetItem from the crawler handler")
		return
	}

	err = data.ToJSON(rw, products)
	if err != nil {
		c.logger.WithError(err).Error("Unable to write getItem using JSON from the crawler handler")
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
		return
	}
	c.logger.WithFields(logrus.Fields{"Handler": c.Identifier, "Function": "getItem"}).Info("Successfully called handler")
}

type KeyItem struct {
}

func (c *Crawler) MiddlewareItemValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		payload, err := data.DecodeToPayload(r)
		if err != nil {
			c.logger.WithError(err).Error("Unable to unmarshal JSON")
			http.Error(rw, "Unable to unmarshal JSON", http.StatusBadRequest)
			return
		}
		if payload.RootURL == "" {
			c.logger.Error("Missing url to crawl. Please set RootURL in payload")
			http.Error(rw, "Missing url to crawl. Please set RootURL in payload", http.StatusBadRequest)
		}
		ctx := context.WithValue(r.Context(), KeyItem{}, payload)
		r = r.WithContext(ctx)

		next.ServeHTTP(rw, r)
	})
}
