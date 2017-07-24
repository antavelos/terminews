package main

import (
	"log"

	"github.com/antavelos/terminews/db"
	c "github.com/jroimartin/gocui"
)

const (
	RSS_READERS_VIEW = "rssreaders"
	NEWS_VIEW        = "news"
	SUMMARY_VIEW     = "summary"
)

var (
	Lists map[string]*List
	tdb   *db.TDB
	// g        *c.Gui
	// err      error
	rrList   *List
	newsList *List
	summary  *c.View
	curW     int
	curH     int
)

func handleFatalError(msg string, err error) {
	// tdb.Close()
	// g.Close()
	log.Fatal(msg, err)
}

func relSize(g *c.Gui) (int, int) {
	tw, th := g.Size()

	return (tw * 3) / 10, (th * 70) / 100

}

// The layout handler calculates all sizes depending
// on the current terminal size.
func layout(g *c.Gui) error {
	// Get the current terminal size.
	tw, th := g.Size()

	// Get the relative size of the views
	rw, rh := relSize(g)

	_, err := g.SetView(RSS_READERS_VIEW, 0, 0, rw, th-1)
	if err != nil {
		handleFatalError("Cannot update rsslist view", err)
	}

	_, err = g.SetView(NEWS_VIEW, rw+1, 0, tw-1, rh)
	if err != nil {
		handleFatalError("Cannot update news view", err)
	}

	_, err = g.SetView(SUMMARY_VIEW, rw+1, rh+1, tw-1, th-1)
	if err != nil {
		handleFatalError("Cannot update summary view.", err)
	}
	updateSummary()

	if curW != tw || curH != th {
		rrList.Draw()
		newsList.Draw()
		curW = tw
		curH = th
	}

	return nil
}

func main() {

	var v *c.View
	var err error

	Lists = make(map[string]*List)

	// Init DB
	if tdb, err = db.InitDB("./term.db"); err != nil {
		tdb.Close()
		log.Fatal(err)
	}
	defer tdb.Close()

	// Create a new GUI.
	g, err := c.NewGui(c.OutputNormal)
	if err != nil {
		handleFatalError("Failed to initialize GUI", err)
	}
	defer g.Close()

	g.SelFgColor = c.ColorGreen
	g.BgColor = c.ColorDefault
	g.Highlight = true

	g.SetManagerFunc(layout)

	curW, curH = g.Size()
	rw, rh := relSize(g)

	// RSS Readers List
	v, err = g.SetView(RSS_READERS_VIEW, 0, 0, rw, curH-1)
	if err != nil && err != c.ErrUnknownView {
		handleFatalError("Failed to create rssreaders list:", err)

	}
	rrList = CreateList(v)

	//
	v, err = g.SetView(NEWS_VIEW, rw+1, 0, curW-1, rh)
	if err != nil && err != c.ErrUnknownView {
		handleFatalError(" Failed to create news list:", err)
	}
	v.Title = " No news yet ... "
	newsList = CreateList(v)

	// Summary view
	summary, err = g.SetView(SUMMARY_VIEW, rw+1, rh+1, curW-1, curH-1)
	if err != nil && err != c.ErrUnknownView {
		handleFatalError("Failed to create summary view:", err)
	}
	summary.Title = " Summary "
	summary.Wrap = true

	loadRssReaders()
	rrList.Focus(g)

	addKeybinding(g, "", rune('a'), c.ModNone, addRssReader)
	addKeybinding(g, "", rune('d'), c.ModNone, deleteEntry)
	addKeybinding(g, "", rune('b'), c.ModNone, addBookmark)
	addKeybinding(g, "", rune('q'), c.ModNone, quit)
	addKeybinding(g, "", c.KeyCtrlB, c.ModNone, loadBookmarks)
	addKeybinding(g, "", c.KeyTab, c.ModNone, switchView)
	addKeybinding(g, "", c.KeyArrowUp, c.ModNone, listUp)
	addKeybinding(g, "", c.KeyArrowDown, c.ModNone, listDown)
	addKeybinding(g, "", c.KeyPgup, c.ModNone, listPgUp)
	addKeybinding(g, "", c.KeyPgdn, c.ModNone, listPgDown)
	addKeybinding(g, "", c.KeyEnter, c.ModNone, loadNews)

	err = g.MainLoop()
	log.Println("terminews exited unexpectedly: ", err)

}
