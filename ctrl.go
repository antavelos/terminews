package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/antavelos/terminews/db"
	c "github.com/jroimartin/gocui"
)

const (
	RSS_READERS_VIEW = "rssreaders"
	NEWS_VIEW        = "news"
	SUMMARY_VIEW     = "summary"
)

var (
	Lists    map[string]*List
	tdb      *db.TDB
	g        *c.Gui
	err      error
	rrList   *List
	newsList *List
	summary  *c.View
)

func handleFatalError(msg string, err error) {
	tdb.Close()
	g.Close()
	log.Fatal(msg, err)
}

func GetList(name string) *List {
	return Lists[name]
}

func CreateViews() {
	tw, th := g.Size()

	lw := (tw * 3) / 10
	oh := (th * 70) / 100

	// RSS Readers List
	rrList, err = CreateList(g, RSS_READERS_VIEW, 0, 0, lw, th-1)
	if err != nil {
		handleFatalError("Failed to create rssreaders list:", err)
	}

	//
	newsList, err = CreateList(g, NEWS_VIEW, lw+1, 0, tw-1, oh)
	if err != nil {
		handleFatalError(" Failed to create news list:", err)
	}
	newsList.Title = " No news yet ... "

	// Summary view
	summary, err = g.SetView(SUMMARY_VIEW, lw+1, oh+1, tw-1, th-1)
	if err != nil && err != c.ErrUnknownView {
		handleFatalError("Failed to create summary view:", err)
	}
	summary.Title = " Summary "
	summary.Wrap = true
}

func UpdateSummary(event db.Event) {
	summary.Clear()
	fmt.Fprintf(summary, "\n\n By %v\n", event.Author)
	fmt.Fprintf(summary, " Published on %v\n\n", event.Published)
	fmt.Fprintf(summary, " %v", event.Summary)
}

func UpdateNews(events []db.Event, from string) {
	data := make([]interface{}, len(events))
	for i, e := range events {
		data[i] = e
	}

	if err = newsList.SetItems(data); err != nil {
		handleFatalError("Failed to update news list", err)
	}
	newsList.SetTitle(fmt.Sprintf("News from %v", from))
	newsList.Focus(g)
	newsList.ResetCursor()
	UpdateSummary(events[0])
}

func LoadRssReaders() []db.RssReader {
	rssReaders, err := tdb.GetRssReaders()
	if err != nil {
		handleFatalError("Failed to load RSS Readers", err)
	}
	data := make([]interface{}, len(rssReaders))
	for i, rr := range rssReaders {
		data[i] = rr
	}

	if err = rrList.SetItems(data); err != nil {
		handleFatalError("Failed to update rss readers list", err)
	}

	rrList.SetTitle("RSS Readers")
	return rssReaders
}

func InitUI() {
	// Create a new GUI.
	g, err = c.NewGui(c.OutputNormal)
	if err != nil {
		handleFatalError("Failed to initialize GUI", err)
	}

	// g.Cursor = true
	g.SelFgColor = c.ColorGreen
	g.BgColor = c.ColorDefault
	g.Highlight = true

	g.SetManagerFunc(layout)
}

func Free() {
	g.Close()
}

func Main() {
	// Init DB
	if tdb, err = db.InitDB("./term.db"); err != nil {
		tdb.Close()
		log.Fatal(err)
	}

	InitUI()
	defer Free()

	CreateViews()

	LoadRssReaders()
	rrList.Focus(g)

	addKeybinding(g, "", rune('a'), c.ModNone, addRssReader)
	addKeybinding(g, "", rune('d'), c.ModNone, deleteRecord)
	addKeybinding(g, "", rune('b'), c.ModNone, bookmark)
	addKeybinding(g, "", rune('q'), c.ModNone, quit)
	addKeybinding(g, "", c.KeyCtrlB, c.ModNone, showBookmarks)
	addKeybinding(g, "", c.KeyTab, c.ModNone, switchView)
	addKeybinding(g, "", c.KeyArrowUp, c.ModNone, listUp)
	addKeybinding(g, "", c.KeyArrowDown, c.ModNone, listDown)
	addKeybinding(g, "", c.KeyPgup, c.ModNone, listPgUp)
	addKeybinding(g, "", c.KeyPgdn, c.ModNone, listPgDown)
	addKeybinding(g, "", c.KeyEnter, c.ModNone, loadNews)
	g.MainLoop()
}

// The layout handler calculates all sizes depending
// on the current terminal size.
func layout(g *c.Gui) error {
	// Get the current terminal size.
	tw, th := g.Size()

	lw := (tw * 3) / 10
	oh := (th * 70) / 100

	_, err := g.SetView(RSS_READERS_VIEW, 0, 0, lw, th-1)
	if err != nil {
		handleFatalError("Cannot update list view", err)
	}
	_, err = g.SetView(NEWS_VIEW, lw+1, 0, tw-1, oh)
	if err != nil {
		handleFatalError("Cannot update output view", err)
	}
	_, err = g.SetView(SUMMARY_VIEW, lw+1, oh+1, tw-1, th-1)
	if err != nil {
		handleFatalError("Cannot update input view.", err)
	}
	return nil
}

func addKeybinding(g *c.Gui, viewname string, key interface{}, mod c.Modifier, handler func(*c.Gui, *c.View) error) {
	err := g.SetKeybinding(viewname, key, mod, handler)
	if err != nil {
		handleFatalError("Could not set key binding:", err)
	}
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

func updateCurrentSummary() {
	currItem := newsList.CurrentItem()
	event := currItem.(db.Event)
	UpdateSummary(event)
}
func listUp(g *c.Gui, v *c.View) error {
	if v == rrList.View {
		rrList.MoveUp()
	} else {
		newsList.MoveUp()
		updateCurrentSummary()
	}

	return nil
}

func listDown(g *c.Gui, v *c.View) error {
	if v == rrList.View {
		rrList.MoveDown()
	} else {
		newsList.MoveDown()
		updateCurrentSummary()
	}

	return nil
}

func listPgDown(g *c.Gui, v *c.View) error {
	if v == rrList.View {
		rrList.MovePgDown()
	} else {
		newsList.MovePgDown()
		updateCurrentSummary()
	}

	return nil
}

func listPgUp(g *c.Gui, v *c.View) error {
	if v == rrList.View {
		rrList.MovePgUp()
	} else {
		newsList.MovePgUp()
		updateCurrentSummary()
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
				UpdateNews(events, rssReader.Name)
			}
			return nil
		})
	}

	return nil
}

func bookmark(g *c.Gui, v *c.View) error {
	if v == newsList.View {
		currItem := newsList.CurrentItem()
		event := currItem.(db.Event)
		if err := tdb.AddEvent(event); err != nil {
			return err
		}
	}
	return nil
}

func showBookmarks(g *c.Gui, v *c.View) error {
	events, err := tdb.GetEvents()
	source := "My bookmarks"
	if err != nil {
		newsList.Title = fmt.Sprintf(" Failed to load news from %v ", source)
		newsList.Clear()
	} else {
		UpdateNews(events, source)
	}
	return nil
}

func addRssReader(g *c.Gui, v *c.View) error {
	return nil
}

func deleteRecord(g *c.Gui, v *c.View) error {
	if v == rrList.View {
		currItem := rrList.CurrentItem()
		rr := currItem.(db.RssReader)
		if err := tdb.DeleteRssReader(rr.Id); err != nil {
			return err
		}
		LoadRssReaders()
	} else {
		if strings.Contains(newsList.Title, "My bookmarks") {
			currItem := newsList.CurrentItem()
			event := currItem.(db.Event)
			if err := tdb.DeleteEvent(event.Id); err != nil {
				return err
			}
			showBookmarks(g, v)
		}
	}
	return nil
}
