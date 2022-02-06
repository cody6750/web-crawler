package handler

import (
	"net/http"

	"github.com/cody6750/web-crawler/web/data"
)

func (c *Crawler) getItem(rw http.ResponseWriter, r *http.Request) error {
	c.logger.Info("Starting to call crawler getItem handler")
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
	c.logger.Info("Successfully called crawler getItem handler")
	return nil
}
