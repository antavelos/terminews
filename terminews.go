package main

import (
	"fmt"
	"github.com/antavelos/terminews/rss"
)

func main() {
	events, _ := rss.Retrieve("http://rss.cnn.com/rss/edition.rss")
	fmt.Println(events)
}
