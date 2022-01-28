package handler

import (
	"log"
	"net/http"

	"github.com/cody6750/codywebapi/webServer/data"
)

// Crawler ...
type Crawler struct {
	logger *log.Logger
}

// NewCrawler ...
func NewCrawler(l *log.Logger) *Crawler {
	return &Crawler{l}
}

// Satisfies HTTP hanlder interface
// Given a request on the http server, based on the path, a specific function is ran. Example  curl -v localhost:9090 will run the function parameter
//http.ResponseWriter : Object that writes to the HTTP response to the user.
//*http.Request : Object containing http request information
func (c *Crawler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet:
		c.getProduct(rw, r)
		return
	default:
		rw.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (c *Crawler) getProduct(rw http.ResponseWriter, r *http.Request) (string, error) {
	products := data.GetProduct()
	err := products.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
	return "", nil
}
