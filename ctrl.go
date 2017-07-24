package main

import (
	"fmt"
	"strings"

	"github.com/antavelos/terminews/db"
	c "github.com/jroimartin/gocui"
)

func addKeybinding(g *c.Gui, viewname string, key interface{}, mod c.Modifier, handler func(*c.Gui, *c.View) error) {
	err := g.SetKeybinding(viewname, key, mod, handler)
	if err != nil {
		fmt.Errorf("Could not set key binding: %v", err)
	}
}

func updateSummary() {
	if newsList.IsEmpty() {
		return
	}
	currItem := newsList.CurrentItem()
	event := currItem.(db.Event)

	summary.Clear()
	fmt.Fprintf(summary, "\n\n By %v\n", event.Author)
	fmt.Fprintf(summary, " Published on %v\n\n", event.Published)
	fmt.Fprintf(summary, " %v", event.Summary)
}

func updateNews(g *c.Gui, events []db.Event, from string) {
	data := make([]interface{}, len(events))
	for i, e := range events {
		data[i] = e
	}

	if err := newsList.SetItems(data); err != nil {
		fmt.Errorf("Failed to update news list: %v", err)
	}
	newsList.SetTitle(fmt.Sprintf("News from %v", from))
	newsList.Focus(g)
	newsList.ResetCursor()
	updateSummary()
}

func loadRssReaders() []db.RssReader {
	rssReaders, err := tdb.GetRssReaders()
	if err != nil {
		fmt.Errorf("Failed to load RSS Readers: %v", err)
	}
	data := make([]interface{}, len(rssReaders))
	for i, rr := range rssReaders {
		data[i] = rr
	}

	if err := rrList.SetItems(data); err != nil {
		fmt.Errorf("Failed to update rss readers list: %v", err)
	}

	rrList.SetTitle("RSS Readers")
	return rssReaders
}

// `quit` is a handler that gets bound to Ctrl-C.
// It signals the main loop to exit.
func quit(g *c.Gui, v *c.View) error {
	return c.ErrQuit
}

func switchView(g *c.Gui, v *c.View) error {
	if v == rrList.View {
		newsList.Focus(g)
		rrList.Unfocus()
	} else {
		rrList.Focus(g)
		newsList.Unfocus()
	}
	return nil
}

func listUp(g *c.Gui, v *c.View) error {
	if v == rrList.View {
		rrList.MoveUp()
	} else {
		newsList.MoveUp()
		updateSummary()
	}

	return nil
}

func listDown(g *c.Gui, v *c.View) error {
	if v == rrList.View {
		rrList.MoveDown()
	} else {
		newsList.MoveDown()
		updateSummary()
	}

	return nil
}

func listPgDown(g *c.Gui, v *c.View) error {
	if v == rrList.View {
		rrList.MovePgDown()
	} else {
		newsList.MovePgDown()
		updateSummary()
	}

	return nil
}

func listPgUp(g *c.Gui, v *c.View) error {
	if v == rrList.View {
		rrList.MovePgUp()
	} else {
		newsList.MovePgUp()
		updateSummary()
	}

	return nil
}

func loadNews(g *c.Gui, v *c.View) error {
	if v == rrList.View {
		currItem := rrList.CurrentItem()
		rssReader := currItem.(db.RssReader)

		newsList.Clear()
		newsList.Focus(g)
		newsList.Title = " Downloading ... "
		g.Execute(func(g *c.Gui) error {
			events, err := DownloadFeed(rssReader.Url)
			if err != nil {
				newsList.Title = fmt.Sprintf(" Failed to load news from %v ", rssReader.Name)
				newsList.Clear()
			} else {
				updateNews(g, events, rssReader.Name)
			}
			return nil
		})
	}

	return nil
}

func addBookmark(g *c.Gui, v *c.View) error {
	if v == newsList.View {
		currItem := newsList.CurrentItem()
		event := currItem.(db.Event)
		if err := tdb.AddEvent(event); err != nil {
			return err
		}
	}
	return nil
}

func loadBookmarks(g *c.Gui, v *c.View) error {
	events, err := tdb.GetEvents()
	source := "My bookmarks"
	if err != nil {
		newsList.Title = fmt.Sprintf(" Failed to load news from %v ", source)
		newsList.Clear()
	} else {
		updateNews(g, events, source)
	}
	return nil
}

func deleteEntry(g *c.Gui, v *c.View) error {
	if v == rrList.View {
		currItem := rrList.CurrentItem()
		rr := currItem.(db.RssReader)
		if err := tdb.DeleteRssReader(rr.Id); err != nil {
			return err
		}
		loadRssReaders()
	} else {
		if strings.Contains(newsList.Title, "My bookmarks") {
			currItem := newsList.CurrentItem()
			event := currItem.(db.Event)
			if err := tdb.DeleteEvent(event.Id); err != nil {
				return err
			}
			loadBookmarks(g, v)
		}
	}
	return nil
}

func addRssReader(g *c.Gui, v *c.View) error {
	return nil
}
