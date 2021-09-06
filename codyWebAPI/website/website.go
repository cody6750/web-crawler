package website

import "log"

//Website ..
type Website interface {
	InitWebsite()
}

func PrintHello() {
	log.Printf("Hello")
}
