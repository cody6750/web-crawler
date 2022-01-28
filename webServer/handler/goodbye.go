package handler

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Goodbye struct {
	log *log.Logger
}

func NewGoodBye(l *log.Logger) *Goodbye {
	return &Goodbye{l}
}

func (g *Goodbye) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	g.log.Print("Goodbye")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, "Bad", 400)
	}
	fmt.Fprintf(rw, "%s\n", body)

}
