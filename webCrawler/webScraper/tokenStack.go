package webcrawler

import (
	"log"

	"golang.org/x/net/html"
)

type stack []html.Token

func (s *stack) push(t html.Token) {
	*s = append(*s, t)
}

func (s *stack) pop() (html.Token, bool) {
	if len(*s) == 0 {
		return html.Token{}, false
	}
	last := len(*s) - 1
	popped := (*s)[last]
	*s = (*s)[:last]
	return popped, true
}

func (s *stack) print() {
	log.Print(*s)
}
