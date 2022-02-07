package handler

import (
	"net/http"

	"github.com/cody6750/web-crawler/web/data"
	"github.com/sirupsen/logrus"
)

func (c *Crawler) getItem(rw http.ResponseWriter, r *http.Request) error {
	c.logger.WithFields(logrus.Fields{"Handler": c.Identifier, "Function": "getItem"}).Info("Starting to call handler")
	payload, err := data.DecodeToPayload(r)
	if err != nil {
		c.logger.WithError(err).Error("Unable to call DecodeToPayload from the crawler handler")
		return err
	}
	products, err := data.GetItem(c.logger, payload.RootURL, payload.ScrapeItemConfiguration, payload.ScrapeURLConfiguration...)
	if err != nil {
		c.logger.WithError(err).Error("Unable to call GetItem from the crawler handler")
		return err
	}
	err = data.ToJSON(rw, products)
	if err != nil {
		c.logger.WithError(err).Error("Unable to write getItem using JSON from the crawler handler")
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
		return err
	}
	c.logger.WithFields(logrus.Fields{"Handler": c.Identifier, "Function": "getItem"}).Info("Successfully called handler")
	return nil
}
