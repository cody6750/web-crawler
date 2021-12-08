package webcrawler

import (
	"encoding/json"
	"log"
)

//Item ...
type Item struct {
	ItemName    string
	URL         string
	TimeQueried string
	DateQueried string
	ItemDetails map[string]string
}

func (i *Item) printJSON() {
	json, _ := json.MarshalIndent(i, "", "    ")
	log.Print("\n" + string(json))
}
