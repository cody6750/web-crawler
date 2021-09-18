package amazon

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	websiteName string = "Amazon"
	WebURL      string = "https://www.amazon.com/s?i=aps&k=RTX%203080&ref=nb_sb_noss_2&url=search-alias%3Daps"
)

//Amazon ... implements the Website Interface
type Amazon struct {
	Name string
}

// Constructor ..
func Constructor() Amazon {
	var amazonObject = Amazon{}
	amazonObject.Name = websiteName
	amazonObject.InitWebsite()
	return amazonObject
}

//InitWebsite ..
func (amazonObject Amazon) InitWebsite() {
	log.Println("Init website")
	amazonObject.Name = websiteName
}

//PrintWebsite ..
func (amazonObject Amazon) PrintWebsite() {
	log.Println("Amazon")
}

//SearchWebsite ..
func (amazonObject Amazon) SearchWebsite(item string) {
	if item == "" {
		log.Fatalf("Item no provided, unable to search website %v", websiteName)
		return
	}
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

	for _, line := range strings.Split(sb, "<") {
		if strings.Contains(line, `span class="a-size-medium a-color-base a-text-normal">`) {
			//if strings.Contains(line, `data-image-index=`) {
			line = strings.TrimPrefix(line, `span class="a-size-medium a-color-base a-text-normal">`)
			fmt.Println(line)
			items = append(items, line)
		}
	}

	//writeToFile(sb)

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
