package main

import (
	"log"
	"net/http"
	"os"

	"github.com/cody6750/codywebapi/webServer/handler"
)

func main() {
	// Given a request on the http server, based on the path, a specific function is ran. Example  curl -v localhost:9090/goodbye will run the function parameter
	l := log.New(os.Stdout, "product-api", log.LstdFlags)
	Hello := handler.NewHello(l)
	sm := http.NewServeMux()
	sm.Handle("/", Hello)
	sm.HandleFunc("/goodbye", func(http.ResponseWriter, *http.Request) {
		log.Print("good bye world")
	})

	//Function for creating a web service on port :9090. Called by curl -v localhost:9090
	http.ListenAndServe(":9090", sm)
	log.Print("done")
}

func tesT() {
	log.Print("hello world")
}
