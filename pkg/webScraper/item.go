package webcrawler

import (
	"encoding/json"
	"log"
	"strings"
)

//Item ...
type Item struct {
	ItemName    string
	URL         *URL
	TimeQueried string
	DateQueried string
	ItemDetails map[string]string
}

func (i *Item) PrintJSON() {
	json, _ := json.MarshalIndent(i, "", "    ")
	output := string(json)
	output = strings.Replace(output, "\\u003c", "<", -1)
	output = strings.Replace(output, "\\u003e", ">", -1)
	output = strings.Replace(output, "\\u0026", "&", -1)
	log.Print("\n" + string(output))
}
