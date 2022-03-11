package data

import (
	"encoding/json"
	"io"

	webcrawler "github.com/cody6750/web-crawler/pkg"
)

// ToJSON encodes the web crawler response to JSON and writes that to the end user over http.
func ToJSON(w io.Writer, r *webcrawler.Response) error {
	e := json.NewEncoder(w)
	e.SetIndent("", "    ")
	return e.Encode(r)
}
