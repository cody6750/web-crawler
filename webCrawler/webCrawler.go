package webcrawler

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	webscraper "github.com/cody6750/codywebapi/webCrawler/webScraper"
)

//WebCrawler ...
type WebCrawler struct {
	client  *http.Client
	scraper *webscraper.WebScraper
}

//New ...
func New() *WebCrawler {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	scraper := webscraper.New()
	crawler := &WebCrawler{
		client:  client,
		scraper: scraper,
	}

	return crawler
}

//Crawl ...
// func (w WebCrawler) Crawl(url string, htmlURLConfig ...webscraper.ExtractURLFromHTMLConfiguration) {
// 	w.scraper.Scrape(url, w.client,
// 		[]webscraper.FormatURLConfiguration{
// 			{},
// 		},
// 		webscraper.ExtractURLFromHTMLConfiguration{
// 			AttributeToCheck:      "class",
// 			AttributeValueToCheck: "a-link-normal a-text-normal",
// 			TagToCheck:            "a",
// 		},
// 		// webscraper.ExtractHTMLLinkConfig{
// 		// 	AttributeToCheck:      "data-ved",
// 		// 	AttributeValueToCheck: "",
// 		// 	TagToCheck:            "a",
// 		// },
// 	)
// }

//Crawl ...
func (w WebCrawler) Crawl(url string, ScrapeConfiguration ...webscraper.ScrapeConfiguration) ([]string, error) {
	// list, _ := w.scraper.Scrape(url, w.client,
	// 	webscraper.ScrapeConfiguration{
	// 		ExtractURLFromHTMLConfiguration: webscraper.ExtractURLFromHTMLConfiguration{
	// 			AttributeToCheck:      "class",
	// 			AttributeValueToCheck: "a-link-normal",
	// 			TagToCheck:            "a",
	// 		},
	// 		FormatURLConfiguration: webscraper.FormatURLConfiguration{
	// 			PrefixToAdd: "http://amazon.com",
	// 		},
	// 	},
	// 	webscraper.ScrapeConfiguration{
	// 		ExtractURLFromHTMLConfiguration: webscraper.ExtractURLFromHTMLConfiguration{
	// 			AttributeToCheck:      "data-routing",
	// 			AttributeValueToCheck: "off",
	// 			TagToCheck:            "a",
	// 		},
	// 		FormatURLConfiguration: webscraper.FormatURLConfiguration{
	// 			PrefixToAdd: "amazon.com",
	// 		},
	// 	},
	// )
	list, _ := w.scraper.Scrape(url, w.client, ScrapeConfiguration...)
	return list, nil
}

//Run ...
func Run() {

	crawl := New()
	crawl.Crawl("https://www.amazon.com/s?k=RTX+3080&ref=nb_sb_noss_2")
	//writeURL("https://www.amazon.com/s?k=RTX+3080&ref=nb_sb_noss_2")
	//crawl.Crawl("https://www.google.com/search?q=RTX+3080&sxsrf=AOaemvJuriGZ27xXjRGoSOpp0evA2muoQw%3A1634776373092&source=hp&ei=NbVwYYK5AtbI1sQPq-qzmAo&iflsig=ALs-wAMAAAAAYXDDRWxhcNlxmfRTq2H2z4aou5VzHdAt&ved=0ahUKEwjCp4LIoNrzAhVWpJUCHSv1DKMQ4dUDCAk&oq=RTX+3080&gs_lcp=Cgdnd3Mtd2l6EAMyBAgjECcyBAgjECcyBAgjECcyCAgAEIAEELEDMggIABCABBCxAzIFCAAQgAQyBQgAEIAEMggIABCABBCxAzILCAAQgAQQsQMQgwEyBQgAEIAEOgcIIxDqAhAnOhEILhCABBCxAxCDARDHARDRAzoOCC4QgAQQsQMQxwEQ0QM6CAgAELEDEIMBUI4XWJ4eYNgeaAFwAHgBgAGhAogB_waSAQUzLjMuMZgBAKABAbABCg&sclient=gws-wiz&uact=5")
	//writeURL("https://www.google.com/search?q=RTX+3080&sxsrf=AOaemvJuriGZ27xXjRGoSOpp0evA2muoQw%3A1634776373092&source=hp&ei=NbVwYYK5AtbI1sQPq-qzmAo&iflsig=ALs-wAMAAAAAYXDDRWxhcNlxmfRTq2H2z4aou5VzHdAt&ved=0ahUKEwjCp4LIoNrzAhVWpJUCHSv1DKMQ4dUDCAk&oq=RTX+3080&gs_lcp=Cgdnd3Mtd2l6EAMyBAgjECcyBAgjECcyBAgjECcyCAgAEIAEELEDMggIABCABBCxAzIFCAAQgAQyBQgAEIAEMggIABCABBCxAzILCAAQgAQQsQMQgwEyBQgAEIAEOgcIIxDqAhAnOhEILhCABBCxAxCDARDHARDRAzoOCC4QgAQQsQMQxwEQ0QM6CAgAELEDEIMBUI4XWJ4eYNgeaAFwAHgBgAGhAogB_waSAQUzLjMuMZgBAKABAbABCg&sclient=gws-wiz&uact=5")
}

func writeURL(url string) {
	resp, err := http.Get(url)
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
