package amazon

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/cody6750/codywebapi/codyWebAPI/tools"
	"golang.org/x/net/html"
)

const (
	WebsiteName                 string = "amazon"
	itemPriceHTMLAttributeValue string = `a-offscreen`
	itemTitleHTMLAttributeValue string = `a-color-base a-text-normal`
	itemURLHTMLAttributeValue   string = `a-link-normal a-text-normal`
	classAttribute              string = `class`
	hrefAttribute               string = `href`
	spanHTMLTag                 string = `span`
	aHTMLTag                    string = `a`
)

//Amazon ... implements the Website Interface
type Amazon struct {
	WebsiteName string
	Name        string
}

type JsonItemOutput struct {
	ItemURL   string `json:"URL"`
	ItemTitle string `json:"Title"`
	ItemPrice string `json:"Price"`
	Date      string `json:"Date"`
}

// New ..
func New() *Amazon {
	amazonObject := &Amazon{}
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
	//searchURL := "https://www.google.com/search?q=RTX+3080&hl=en&sxsrf=AOaemvKDvMp_Dp95q3Mbd5I-f8xfuHHlaQ%3A1634414248578&source=hp&ei=qC5rYcraIOCXr7wP5q-VmAk&iflsig=ALs-wAMAAAAAYWs8uAtN0F-7iIy_fsMottIa6a9orwQh&ved=0ahUKEwjKzszF28_zAhXgy4sBHeZXBZMQ4dUDCAk&uact=5&oq=RTX+3080&gs_lcp=Cgdnd3Mtd2l6EAMyBAgjECcyBAgjECcyBAgjECcyBwguELEDEEMyDQgAEIAEEIcCELEDEBQyCAgAEIAEELEDMgQIABBDMgUIABCABDIFCAAQgAQyCAgAEIAEELEDOgQILhBDOg0ILhCxAxDHARDRAxBDOgcIABCxAxBDOgoIABCxAxDJAxBDOgoIABCABBCHAhAUUI39O1jPhDxgtIU8aABwAHgAgAGgA4gB8AuSAQcyLTMuMS4xmAEAoAEB&sclient=gws-wiz"
	log.Printf("%v successfully generated search URL %v", tools.FuncName(), searchURL)
	return searchURL, nil
}

func generateItemURLFromHref(hrefValue string) (string, error) {
	if hrefValue == "" {
		log.Printf("%v unable to call function, item is empty", tools.FuncName())
		return hrefValue, nil
	}
	itemURL := "https://www.amazon.com" + hrefValue
	//log.Printf("%v successfully generated search URL %v", tools.FuncName(), itemURL)
	return itemURL, nil
}

func getHTTPAttributeValueFromToken(t html.Token, attributeValueToGet string) (attributeValue string, err error) {
	for _, a := range t.Attr {
		if a.Key == attributeValueToGet {
			attributeValue = a.Val
		}
	}
	if attributeValue == "" {
		return attributeValue, errors.New("TODO")
	}
	return attributeValue, nil
}

//SearchWebsite ..
func (amazonObject Amazon) SearchWebsite(item string) ([]string, error) {
	var (
		itemURLCheck                                  map[string]bool
		items                                         []string
		attributeValue, itemURL, itemPrice, itemTitle string
		correctText, search                           bool = false, true
	)
	itemURLCheck = make(map[string]bool)
	if item == "" {
		log.Printf("Item no provided, unable to search website %v", WebsiteName)
		return items, nil
	}
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	WebURL, _ := generateSearchURL(item)
	// Create and modify HTTP request before sending
	request, err := http.NewRequest("GET", WebURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("User-Agent", "This bot just searches amazon for a product")

	// Make request
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	z := html.NewTokenizer(response.Body)
	for search {
		tt := z.Next()
		switch {
		case tt == html.TextToken && correctText:
			if strings.Contains(attributeValue, itemTitleHTMLAttributeValue) {
				itemTitle = string(z.Text())
			} else if strings.Contains(attributeValue, itemPriceHTMLAttributeValue) {
				itemPrice = string(z.Text())
			}
			//log.Print(string(z.Text()))
		case tt == html.ErrorToken:
			search = false
		case tt == html.StartTagToken:
			t := z.Token()
			if t.Data == spanHTMLTag {
				attributeValue, err = getHTTPAttributeValueFromToken(t, classAttribute)
				if err != nil {
					break
				}
				// Used to obtain item title and price from the span tags, using class attirbute values.
				if strings.Contains(attributeValue, itemTitleHTMLAttributeValue) || strings.Contains(attributeValue, itemPriceHTMLAttributeValue) {
					correctText = true
				}
			} else if t.Data == aHTMLTag {
				attributeValue, err = getHTTPAttributeValueFromToken(t, classAttribute)
				if strings.Contains(attributeValue, itemURLHTMLAttributeValue) {
					attributeValue, err = getHTTPAttributeValueFromToken(t, hrefAttribute)
					if err != nil {
						break
					}
					itemURL, err = generateItemURLFromHref(attributeValue)
					if err != nil {
						break
					}
					if _, exist := itemURLCheck[itemURL]; exist {
						continue
					}
					itemURLCheck[itemURL] = true
				}
			}
		case tt == html.EndTagToken:
			t := z.Token()
			if t.Data == spanHTMLTag || t.Data == aHTMLTag {
				correctText = false
			}
		}
		if itemURL != "" && itemPrice != "" && itemTitle != "" {
			log.Print("Generating json")
			itemJSONOutput := &JsonItemOutput{
				ItemURL:   itemURL,
				ItemPrice: itemPrice,
				ItemTitle: itemTitle,
				Date:      time.Now().String(),
			}

			b := new(bytes.Buffer)
			encoder := json.NewEncoder(b)
			encoder.SetEscapeHTML(false)
			encoder.Encode(itemJSONOutput)
			log.Print(b)
			itemURL, itemPrice, itemTitle = "", "", ""
		}

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
	writeToFile(sb)
	return items, nil
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
