package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/cody6750/codywebapi/webServer/handler"
)

func main() {
	// Given a request on the http server, based on the path, a specific function is ran. Example  curl -v localhost:9090/goodbye will run the function parameter
	l := log.New(os.Stdout, "product-api", log.LstdFlags)
	Crawler := handler.NewCrawler(l)
	sm := http.NewServeMux()
	sm.Handle("/", Crawler)

	//Function for creating a web service on port :9090. Called by curl -v localhost:9090
	server := http.Server{
		Addr:         ":9090",
		Handler:      sm,
		IdleTimeout:  time.Second * 120,
		ReadTimeout:  time.Second * 1,
		WriteTimeout: time.Second * 1,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)
	sig := <-sigChan
	l.Print("Recieved teriminate, graceful shutdown", sig)
	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	server.Shutdown(tc)

}
