package main

import (
	"log"
)

func main() {
	log.Println("hi")
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
