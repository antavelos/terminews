package rss

import (
	"errors"
	"fmt"
	"github.com/antavelos/terminews/news"
	"github.com/mmcdole/gofeed"
)

func Retrieve(url string) (news.Events, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Couldn't retrieve data from: '%v'", url))
	}

	var events news.Events
	for _, item := range feed.Items {
		e := news.Event{
			item.Title,
			item.Link,
			item.Description,
			item.Published,
		}
		events = append(events, e)
	}

	return events, nil
}
