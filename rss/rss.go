package rss

import (
	"errors"
	"fmt"
	"github.com/antavelos/terminews/db"
	"github.com/mmcdole/gofeed"
	"regexp"
)

func Retrieve(url string) ([]db.Event, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
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
		e.Summary = trimSummary(item.Description)
		e.Published = item.Published

		events = append(events, e)
	}

	return events, nil
}

func trimSummary(desc string) string {
	var re = regexp.MustCompile(`(<.*?>)`)

	return re.ReplaceAllString(desc, ``)
}
