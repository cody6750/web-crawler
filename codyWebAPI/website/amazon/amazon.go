package amazon

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/cody6750/codywebapi/codyWebAPI/tools"
)

const (
	WebsiteName string = "amazon"
)

//Amazon ... implements the Website Interface
type Amazon struct {
	WebsiteName string
	Name        string
}

// Constructor ..
func Constructor() Amazon {
	var amazonObject = Amazon{}
	amazonObject.Name = WebsiteName
	amazonObject.InitWebsite()
	return amazonObject
}

//InitWebsite ..
func (amazonObject Amazon) InitWebsite() {
	log.Println("Init website")
	amazonObject.Name = WebsiteName
}

//PrintWebsite ..
func (amazonObject Amazon) PrintWebsite() {
	log.Println("Amazon")
}

func generateSearchURL(item string) (string, error) {
	if item == "" {
		log.Printf("%v unable to call function, item is empty", tools.FuncName())
		return item, nil
	}
	searchURL := "https://www.amazon.com/s?k=" + strings.ReplaceAll(item, " ", "+") + "&ref=nb_sb_noss_2"
	log.Printf("%v successfully generated search URL %v", tools.FuncName(), searchURL)
	return searchURL, nil
}

//SearchWebsite ..
func (amazonObject Amazon) SearchWebsite(item string) {
	if item == "" {
		log.Printf("Item no provided, unable to search website %v", WebsiteName)
		return
	}
	WebURL, _ := generateSearchURL(item)
	resp, err := http.Get(WebURL)
	if err != nil {
		log.Fatalln(err)
	}

	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	//Convert the body to type string
	sb := string(body)
	var items []string
	var size string
	for _, line := range strings.Split(sb, "<") {
		//log.Printf(line)
		if strings.Contains(line, `span class="a-size-medium a-color-base a-text-normal">`) || strings.Contains(line, `<span class="a-size-base-plus a-color-base a-text-normal">`) {
			if strings.Contains(line, `span class="a-size-medium a-color-base a-text-normal">`) {
				size = "medium"
			}
			if strings.Contains(line, `<span class="a-size-base-plus a-color-base a-text-normal">`) {
				size = "base"
			}

			correctHTML := `span class="a-size-` + size + `medium a-color-base a-text-normal">`
			line = strings.TrimPrefix(line, correctHTML)
			log.Print(line)
			items = append(items, line)
		}
	}

	writeToFile(sb)

}

func writeToFile(body string) {
	f, err := os.Create("data.html")

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	_, err2 := f.WriteString(body)

	if err2 != nil {
		log.Fatal(err2)
	}

	fmt.Println("done")
}

func helpFunction() {
	log.Printf("Listing all avaliable commands:")
}
