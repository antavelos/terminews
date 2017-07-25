package main

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/antavelos/terminews/db"
	"github.com/mmcdole/gofeed"
)

func CheckUrl(url string) (*gofeed.Feed, error) {
	fp := gofeed.NewParser()
	return fp.ParseURL(url)
	// if err != nil {
	// 	return nil, false
	// }
	// return feed, true
}

func DownloadEvents(url string) ([]db.Event, error) {
	feed, err := CheckUrl(url)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to retrieve news from: '%v'", url))
	}

	var events []db.Event
	for _, item := range feed.Items {
		e := db.Event{}
		e.Title = item.Title
		if item.Author != nil {
			e.Author = item.Author.Name
		} else {
			e.Author = "Unknown"
		}
		e.Url = item.Link
		e.Summary = trimHtml(item.Description)
		e.Published = item.Published

		events = append(events, e)
	}

	return events, nil
}

func trimHtml(desc string) string {
	var re = regexp.MustCompile(`(<.*?>)`)

	return re.ReplaceAllString(desc, ``)
}
