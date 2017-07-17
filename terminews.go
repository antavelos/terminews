package main

import (
	// "flag"
	// "fmt"
	// "github.com/antavelos/terminews/rss"
	"github.com/antavelos/terminews/ui"
	// "github.com/fatih/color"
)

func main() {
	ui.Ui()
	// flag.Usage = func() {
	// 	fmt.Println("foo bar")
	// }
	// c := flag.String("f", "default", "usage")
	// flag.Parse()
	// fmt.Println(*c)
	// events, err := rss.Retrieve("http://rss.cnn.com/rss/edition.rss")
	// if err != nil {
	// 	red := color.New(color.FgRed, color.Bold).SprintFunc()

	// 	fmt.Print(red(err))
	// } else {
	// 	green := color.New(color.FgGreen, color.Bold).SprintFunc()
	// 	fmt.Println(events)
	// 	found := fmt.Sprintf("\nFound %v news", len(events))
	// 	fmt.Println(green(found))
	// }
}
