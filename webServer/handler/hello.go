package handler

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Hello struct {
	log *log.Logger
}

func NewHello(l *log.Logger) *Hello {
	return &Hello{l}
}

// Satisfies HTTP hanlder interface
func (h *Hello) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// Given a request on the http server, based on the path, a specific function is ran. Example  curl -v localhost:9090 will run the function parameter
	//http.ResponseWriter : Object that writes to the HTTP response to the user.
	//*http.Request : Object containing http request information
	h.log.Print("Hello")
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		//Allows you to return specific error codes
		http.Error(rw, "bad request", 400)
		return
	}
	fmt.Fprintf(rw, "%s\n", b)
}
