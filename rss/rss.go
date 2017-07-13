package rss

import (
	"github.com/antavelos/terminews/news"
	"github.com/mmcdole/gofeed"
)

func Retrieve(url string) (news.Events, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		return nil, err
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
