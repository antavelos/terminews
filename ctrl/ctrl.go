package ctrl

import (
	"bytes"
	"fmt"
	"log"

	"github.com/antavelos/terminews/db"
	"github.com/antavelos/terminews/news"
	"github.com/antavelos/terminews/rss"
	"github.com/antavelos/terminews/ui"
	c "github.com/jroimartin/gocui"
)

const (
	RSS_READERS_VIEW = "rssreaders"
	NEWS_VIEW        = "news"
	SUMMARY_VIEW     = "summary"
)

// Items to fill the list with.
var listItems = []string{
	"Line 1",
	"Line 2",
	"Line 3",
	"Line 4",
	"Line 5",
}

var (
	tdb      *db.TDB
	g        *c.Gui
	err      error
	rrList   *ui.List
	newsList *ui.List
	summary  *c.View
)

func handleFatalError(msg string, err error) {
	tdb.Close()
	g.Close()
	log.Fatal(msg, err)
}

func spaces(n int) string {
	var s bytes.Buffer
	for i := 0; i < n; i++ {
		s.WriteString(" ")
	}
	return s.String()

}

func CreateViews() {
	tw, th := g.Size()

	lw := (tw * 3) / 10
	oh := (th * 6) / 10

	// RSS Readers List
	rrList, err = ui.CreateList(g, RSS_READERS_VIEW, 0, 0, lw, th-1)
	if err != nil {
		handleFatalError("Failed to create rssreaders list:", err)
	}
	rrList.Title = " RSS Readers "

	//
	newsList, err = ui.CreateList(g, NEWS_VIEW, lw+1, 0, tw-1, oh)
	if err != nil {
		handleFatalError(" Failed to create news list:", err)
	}
	newsList.Title = "News from ..."

	// Summary view
	summary, err = g.SetView(SUMMARY_VIEW, lw+1, oh+1, tw-1, th-1)
	if err != nil && err != c.ErrUnknownView {
		handleFatalError("Failed to create summary view:", err)
	}
	summary.Title = " Summary "
	summary.Wrap = true
}

func UpdateSummary(event news.Event) {
	summary.Clear()
	fmt.Fprintf(summary, "\n\n By %v\n", event.Author)
	fmt.Fprintf(summary, " Published on %v\n\n", event.Published)
	fmt.Fprintf(summary, " %v", event.Description)
}

func UpdateNews(rr db.RssReader) {
	events, err := rss.Retrieve(rr.Url)
	if err != nil {
		newsList.Title = fmt.Sprintf(" Failed to load news from %v ", rr.Name)
		newsList.Clear()
	}
	var data []ui.Displayer = make([]ui.Displayer, len(events))
	for i, e := range events {
		data[i] = e
	}

	if err = newsList.SetItems(data); err != nil {
		handleFatalError("Failed to update news list", err)
	}
	newsList.Title = fmt.Sprintf(" News from %v ", rr.Name)
	newsList.Focus(g)
	UpdateSummary(events[0])
}

func LoadRssReaders() []db.RssReader {
	rssReaders, err := tdb.GetRssReaders()
	if err != nil {
		handleFatalError("Failed to load RSS Readers", err)
	}
	var data []ui.Displayer = make([]ui.Displayer, len(rssReaders))
	for i, rr := range rssReaders {
		data[i] = rr
	}

	if err = rrList.SetItems(data); err != nil {
		handleFatalError("Failed to update rss readers list", err)
	}

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

// Set up the widgets and run the event loop.
func Main() {
	// Init DB
	if tdb, err = db.InitDB("./terminews.db"); err != nil {
		tdb.Close()
		log.Fatal(err)
	}

	InitUI()
	defer Free()

	CreateViews()

	rssReaders := LoadRssReaders()

	UpdateNews(rssReaders[0])

	err = g.SetKeybinding("", rune('q'), c.ModNone, quit)
	if err != nil {
		handleFatalError("Could not set key binding:", err)
	}

	err = g.SetKeybinding("", c.KeyTab, c.ModNone, switchView)
	if err != nil {
		handleFatalError("Could not set key binding:", err)
	}

	err = g.SetKeybinding("", c.KeyArrowUp, c.ModNone, listUp)
	if err != nil {
		handleFatalError("Could not set key binding:", err)
	}

	err = g.SetKeybinding("", c.KeyArrowDown, c.ModNone, listDown)
	if err != nil {
		handleFatalError("Could not set key binding:", err)
	}

	err = g.SetKeybinding("", c.KeyEnter, c.ModNone, loadNews)
	if err != nil {
		handleFatalError("Could not set key binding:", err)
	}
	// Start the main loop.
	g.MainLoop()
}

// The layout handler calculates all sizes depending
// on the current terminal size.
func layout(g *c.Gui) error {
	// Get the current terminal size.
	tw, th := g.Size()

	lw := (tw * 3) / 10
	oh := (th * 6) / 10

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
		currItem := newsList.CurrentItem()
		event := currItem.(news.Event)
		UpdateSummary(event)
	}

	return nil
}

func listDown(g *c.Gui, v *c.View) error {
	if v == rrList.View {
		rrList.MoveDown()
	} else {
		newsList.MoveDown()
		currItem := newsList.CurrentItem()
		event := currItem.(news.Event)
		UpdateSummary(event)
	}

	return nil
}

func loadNews(g *c.Gui, v *c.View) error {

	if v == rrList.View {
		currItem := rrList.CurrentItem()
		rssReader := currItem.(db.RssReader)
		UpdateNews(rssReader)
	}

	return nil
}
