package webcrawler

import (
	"net/http"
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
func (w WebCrawler) Crawl(url string) {
	w.scraper.Scrape(url, w.client, webscraper.ExtractHtmlLinkConfig{
		AttributeToCheck:      "class",
		AttributeValueToCheck: "a-size-medium",
		TagToCheck:            "span",
	},
	// webscraper.ExtractHtmlLinkConfig{
	// 	AttributeToCheck:      "class",
	// 	AttributeValueToCheck: "nav-a",
	// 	TagToCheck:            "a",
	// },
	)
}

func Main() {
	crawl := New()
	crawl.Crawl("https://www.amazon.com/s?k=RTX+3080&ref=nb_sb_noss_2")

	//crawl.Crawl("https://www.google.com/search?q=RTX+3080&sxsrf=AOaemvJuriGZ27xXjRGoSOpp0evA2muoQw%3A1634776373092&source=hp&ei=NbVwYYK5AtbI1sQPq-qzmAo&iflsig=ALs-wAMAAAAAYXDDRWxhcNlxmfRTq2H2z4aou5VzHdAt&ved=0ahUKEwjCp4LIoNrzAhVWpJUCHSv1DKMQ4dUDCAk&oq=RTX+3080&gs_lcp=Cgdnd3Mtd2l6EAMyBAgjECcyBAgjECcyBAgjECcyCAgAEIAEELEDMggIABCABBCxAzIFCAAQgAQyBQgAEIAEMggIABCABBCxAzILCAAQgAQQsQMQgwEyBQgAEIAEOgcIIxDqAhAnOhEILhCABBCxAxCDARDHARDRAzoOCC4QgAQQsQMQxwEQ0QM6CAgAELEDEIMBUI4XWJ4eYNgeaAFwAHgBgAGhAogB_waSAQUzLjMuMZgBAKABAbABCg&sclient=gws-wiz&uact=5")
}
