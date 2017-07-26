package main

import (
	"fmt"
	"strings"
	_ "time"

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
	fmt.Fprintf(summary, "\n\n %v %v\n", Bold.Sprint("By"), event.Author)
	fmt.Fprintf(summary, " %v %v\n\n", Bold.Sprint("Published on"), event.Published)
	fmt.Fprintf(summary, " %v", event.Summary)
}

func updateNews(g *c.Gui, events []db.Event, from string) {
	if len(events) == 0 {
		newsList.SetTitle(fmt.Sprintf("No news in %v", from))
		summary.Clear()
		return
	}
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

func createPromptView(g *c.Gui, title string) error {
	tw, th := g.Size()
	v, err := g.SetView(PROMPT_VIEW, tw/6, (th/2)-1, (tw*5)/6, (th/2)+1)
	if err != nil && err != c.ErrUnknownView {
		return err
	}
	v.Editable = true
	setPromptViewTitle(g, title)
	// v.Title = title

	g.Cursor = true
	g.SetCurrentView(PROMPT_VIEW)

	return nil
}

func deletePromptView(g *c.Gui) error {
	if err := g.DeleteView(PROMPT_VIEW); err != nil {
		return err
	}
	g.Cursor = false

	return nil
}

func setPromptViewTitle(g *c.Gui, title string) {
	v, _ := g.View(PROMPT_VIEW)
	v.Title = fmt.Sprintf("%v (Ctrl-q to cancel)", title)
}

func isNewRSSPrompt(v *c.View) bool {
	return strings.Contains(v.Title, "new RSS")
}

func isFindPrompt(v *c.View) bool {
	return strings.Contains(v.Title, "Search ")
}

func termsInEvent(terms []string, e db.Event) bool {
	for _, t := range terms {
		if strings.Contains(e.Title, t) || strings.Contains(e.Summary, t) {
			return true
		}
	}
	return false
}

func findEvents(terms []string) chan db.Event {
	c := make(chan db.Event)
	readers, _ := tdb.GetRssReaders()
	// if err != nil {
	// 	return err
	// }
	go func() {
		for _, reader := range readers {
			events, err := DownloadEvents(reader.Url)
			if err != nil {
				continue
			}
			for _, e := range events {
				if termsInEvent(terms, e) {
					c <- e
				}
			}
		}
		close(c)
	}()
	return c
}

// Key binding functions

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
		if !newsList.IsEmpty() {
			newsList.MoveUp()
			updateSummary()
		}
	}

	return nil
}

func listDown(g *c.Gui, v *c.View) error {
	if v == rrList.View {
		rrList.MoveDown()
	} else {
		if !newsList.IsEmpty() {
			newsList.MoveDown()
			updateSummary()
		}
	}

	return nil
}

func listPgDown(g *c.Gui, v *c.View) error {
	if v == rrList.View {
		rrList.MovePgDown()
	} else {
		if !newsList.IsEmpty() {
			newsList.MovePgDown()
			updateSummary()
		}
	}

	return nil
}

func listPgUp(g *c.Gui, v *c.View) error {
	if v == rrList.View {
		rrList.MovePgUp()
	} else {
		if !newsList.IsEmpty() {
			newsList.MovePgUp()
			updateSummary()
		}
	}

	return nil
}

func onEnter(g *c.Gui, v *c.View) error {
	switch v.Name() {
	case RSS_READERS_VIEW:
		currItem := rrList.CurrentItem()
		rssReader := currItem.(db.RssReader)

		newsList.Clear()
		newsList.Focus(g)
		newsList.Title = " Downloading ... "
		g.Execute(func(g *c.Gui) error {
			events, err := DownloadEvents(rssReader.Url)
			if err != nil {
				newsList.Title = fmt.Sprintf(" Failed to load news from %v ", rssReader.Name)
				newsList.Clear()
			} else {
				updateNews(g, events, rssReader.Name)
			}
			return nil
		})
	case PROMPT_VIEW:
		if isNewRSSPrompt(v) {
			url := strings.TrimSpace(v.ViewBuffer())
			if len(url) == 0 {
				return nil
			}
			g.Execute(func(g *c.Gui) error {
				feed, err := CheckUrl(url)
				if err != nil {
					setPromptViewTitle(g, "Invalid URL, try again:")
					g.SelFgColor = c.ColorRed | c.AttrBold
					return nil
				}

				_, err = tdb.GetRssReaderByUrl(url)
				if _, ok := err.(db.NotFound); !ok {
					setPromptViewTitle(g, "RSS Reader already exists, try again:")
					g.SelFgColor = c.ColorRed | c.AttrBold
					return nil
				}

				rr := db.RssReader{Name: feed.Title, Url: url}
				if err := tdb.AddRssReader(rr); err != nil {
					return err
				}
				deletePromptView(g)
				g.SelFgColor = c.ColorGreen | c.AttrBold
				loadRssReaders()
				rrList.Focus(g)

				return nil
			})
		}
		if isFindPrompt(v) {
			newsList.Reset()
			newsList.Focus(g)
			newsList.Title = " Searching ... "
			deletePromptView(g)
			terms := strings.Split(strings.TrimSpace(v.ViewBuffer()), " ")
			g.Execute(func(g *c.Gui) error {
				c := 0
				for event := range findEvents(terms) {
					newsList.AddItem(g, event)
					c++
				}
				if c == 0 {
					newsList.SetTitle("No events found")
				} else {
					newsList.SetTitle(fmt.Sprintf("%v events found", c))
				}
				return nil
			})
		}
	}

	return nil
}

func addBookmark(g *c.Gui, v *c.View) error {
	switch v.Name() {
	case NEWS_VIEW:
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

func removePrompt(g *c.Gui, v *c.View) error {
	if v.Name() == PROMPT_VIEW {
		rrList.Focus(g)
		g.SelFgColor = c.ColorGreen | c.AttrBold
		return deletePromptView(g)
	}
	return nil
}

func addRssReader(g *c.Gui, v *c.View) error {
	if err := createPromptView(g, "Give a new RSS reader URL:"); err != nil {
		return err
	}

	return nil
}

func find(g *c.Gui, v *c.View) error {
	if err := createPromptView(g, "Search with multiple terms:"); err != nil {
		return err
	}

	return nil
}
