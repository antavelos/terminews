package rss

import (
	"errors"
	"fmt"
	"github.com/antavelos/terminews/news"
	"github.com/mmcdole/gofeed"
	"regexp"
)

func Retrieve(url string) ([]news.Event, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to retrieve news from: '%v'", url))
	}

	var events []news.Event
	for _, item := range feed.Items {
		e := news.Event{}
		e.Title = removeSpecialCharacters(item.Title)
		if item.Author != nil {
			e.Author = removeSpecialCharacters(item.Author.Name)
		} else {
			e.Author = "Unknown"
		}
		e.Link = item.Link
		e.Description = trimDescription(removeSpecialCharacters(item.Description))
		e.Published = removeSpecialCharacters(item.Published)

		events = append(events, e)
	}

	return events, nil
}

func trimDescription(desc string) string {
	var re = regexp.MustCompile(`(<.*?>)`)

	return re.ReplaceAllString(desc, ``)
}

func removeSpecialCharacters(s string) string {
	var re = regexp.MustCompile(`([’‘])`)

	return re.ReplaceAllString(s, `'`)
}
