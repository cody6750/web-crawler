package main

import "github.com/cody6750/web-crawler/web/server"

func main() {
	webCrawlerServer := server.New()
	webCrawlerServer.Run()
}
