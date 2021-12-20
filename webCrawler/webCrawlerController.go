package webcrawler

import (
	"log"
	"strconv"

	webscraper "github.com/cody6750/codywebapi/webCrawler/webScraper"
)

//Queue ...
type Queue []string

func (q *Queue) enqueue(s string) {
	*q = append(*q, s)
}

func (q *Queue) enqueueList(s []string) {
	for _, item := range s {
		*q = append(*q, item)
	}
}
func (q *Queue) dequeue() string {
	if len(*q) == 0 {
		return ""
	}
	first := (*q)[0]
	*q = (*q)[1:]
	return first
}

func (q *Queue) print() {
	log.Print(*q)
}

func testQueue() {
	var q Queue
	s := []string{
		"he111llo",
		"h1111i",
	}
	q.enqueueList(s)
	q.print()
	for i := 0; i < 6; i++ {
		q.enqueue("hello" + strconv.Itoa(i))
	}
	q.print()
	q.enqueueList(s)
	q.print()
	log.Print(q.dequeue())
	q.print()
	q.enqueue("hello" + strconv.Itoa(7))
	q.print()

}

//WebCrawlController ...
func WebCrawlController() {
	// var counter int
	duplicateUrls := make(chan map[string]bool)
	// quit := make(chan int, 1)
	// var wg sync.WaitGroup
	crawler := New()
	crawler.Crawl("https://www.bestbuy.com/site/searchpage.jsp?st=RTX+3080&_dyncharset=UTF-8&_dynSessConf=&id=pcat17071&type=page&sc=Global&cp=1&nrp=&sp=&qp=&list=n&af=true&iht=y&usc=All+Categories&ks=960&keys=keys",
		duplicateUrls,
		2,
		[]webscraper.ScrapeItemConfiguration{
			{
				ItemName: "Graphics Cards",
				ItemToGet: webscraper.ExtractFromHTMLConfiguration{
					Tag:            "li",
					Attribute:      "class",
					AttributeValue: "sku-item",
				},
				ItemDetails: map[string]webscraper.ExtractFromHTMLConfiguration{
					"title": {
						Tag:            "h4",
						Attribute:      "class",
						AttributeValue: "sku-header",
					},
					"price": {
						Tag:            "span",
						Attribute:      "aria-hidden",
						AttributeValue: "true",
					},
					"link": {
						Tag:            "a",
						Attribute:      "",
						AttributeValue: "",
						AttributeToGet: "href",
					},
					"In stock": {
						Tag:            "button",
						AttributeToGet: "data-button-state",
						AttributeValue: "button",
						Attribute:      "disabled type",
					},
					"Out of stock": {
						Tag:            "button",
						AttributeToGet: "data-button-state",
						AttributeValue: "button",
						Attribute:      "type",
					},
				},
			},
		},
		[]webscraper.ScrapeURLConfiguration{
			{
				// ExtractFromHTMLConfiguration: ExtractFromHTMLConfiguration{
				// 	Attribute:      "class",
				// 	AttributeValue: "a-link-normal",
				// 	Tag:            "a",
				// },
				FormatURLConfiguration: webscraper.FormatURLConfiguration{
					PrefixExist: "/",
					PrefixToAdd: "http://bestbuy.com",
				},
			},
		}...,
	)
	// queueToPass.enqueueList(list)
	// urlTocrawl <- queueToPass
	// for {
	// 	select {
	// 	case currentQueue := <-urlTocrawl:
	// 		wg.Add(1)
	// 		go func() {
	// 			counter++
	// 			defer wg.Done()
	// 			if len(currentQueue) == 0 {
	// 				quit <- 1
	// 			}
	// 			log.Print(currentQueue.dequeue())
	// 			urlTocrawl <- currentQueue
	// 		}()
	// 	case <-quit:
	// 		wg.Wait()
	// 		log.Printf("Pod count :%v", counter)
	// 		close(urlTocrawl)
	// 		close(quit)
	// 		return
	// 	}
	// }
}
