package ctrl

import (
	"fmt"
	"github.com/antavelos/terminews/db"
	"github.com/antavelos/terminews/news"
	"github.com/antavelos/terminews/rss"
	"github.com/antavelos/terminews/ui"
	gc "github.com/rthornton128/goncurses"
	"log"
)

const (
	RSS_READERS_WINDOW = iota
	NEWS_WINDOW
)

const (
	FOCUSED_WINDOW   = 1
	UNFOCUSED_WINDOW = 2
	ERROR_WINDOW     = 3
)

var (
	tdb                                *db.TDB
	stdscr                             *gc.Window
	rssReadersWin, newsWin             *ui.TWindow
	rssReadersMenu, newsMenu           *ui.TMenu
	rssReadersMenuItems, newsMenuItems []*ui.TMenuItem
	windows                            []*ui.TWindow
	panels                             []*gc.Panel
	activeWindow                       int
	err                                error
	maxX                               int
	maxY                               int
)

func CreaterssReadersWindow(rssReaders []db.RssReader) error {

	// MenuItems
	for i, rssReader := range rssReaders {
		item := &ui.TMenuItem{}
		if err = item.Create(i+1, &rssReader); err != nil {
			return err
		}
		rssReadersMenuItems = append(rssReadersMenuItems, item)
	}
	// Menu
	rssReadersMenu = &ui.TMenu{}
	if err = rssReadersMenu.Create(rssReadersMenuItems); err != nil {
		return err
	}
	// Window
	rssReadersWin = &ui.TWindow{}
	if err = rssReadersWin.Create("My RSS Readers", maxY, maxX/4, 0, 0); err != nil {
		return err
	}
	windows = append(windows, rssReadersWin)

	// attach rssReadersMenu on rssReadersWindow
	rssReadersWin.AttachMenu(rssReadersMenu)

	panels = append(panels, gc.NewPanel(rssReadersWin.Window))

	rssReadersMenu.Post()
	rssReadersWin.Refresh()
	rssReadersWin.Focus(FOCUSED_WINDOW)

	return nil
}

func CreateNewsData(events []news.Event, rrName string) error {
	// MenuItems
	for i, event := range events {
		item := &ui.TMenuItem{}
		if err = item.Create(i+1, &event); err != nil {
			return err
		}
		newsMenuItems = append(newsMenuItems, item)
	}
	// Menu
	newsMenu = &ui.TMenu{}
	if err = newsMenu.Create(newsMenuItems); err != nil {
		return err
	}

	newsWin = &ui.TWindow{}
	title := fmt.Sprintf("News from %v", rrName)
	if err = newsWin.Create(title, maxY, (maxX*3)/4, 0, maxX/4); err != nil {
		return err
	}
	windows = append(windows, newsWin)

	// attach newsMenu on newsWindow
	newsWin.AttachMenu(newsMenu)

	panels = append(panels, gc.NewPanel(newsWin.Window))

	newsMenu.Post()
	newsWin.Refresh()

	return nil
}

// the order mstters!!!
func Free() {
	var tmi *ui.TMenuItem

	rssReadersMenu.UnPost()
	newsMenu.UnPost()

	for _, tmi = range rssReadersMenuItems {
		tmi.MenuItem.Free()
	}
	for _, tmi = range newsMenuItems {
		tmi.MenuItem.Free()
	}

	rssReadersMenu.Free()
	newsMenu.Free()

	gc.End()
}

func handleUIFatalError(err error) {
	gc.End()
	log.Fatal(err)
}

func handleDBFatalError(err error) {
	tdb.Close()
	log.Fatal(err)
}

func handleNewsLoadError(rrName string) {
	msg := fmt.Sprintf("Failed to load news from %v", rrName)
	newsWin.SetTitle(msg)
	newsWin.Focus(ERROR_WINDOW)
	activeWindow = NEWS_WINDOW
}

func InitUI() error {

	stdscr, err = gc.Init()
	if err != nil {
		return err
	}

	maxY, maxX = stdscr.MaxYX()
	// Initial configuration
	gc.StartColor()
	gc.Raw(true)
	gc.Echo(false)
	gc.Cursor(0)
	stdscr.Keypad(true)
	stdscr.Clear()

	// Define color combinations
	gc.InitPair(FOCUSED_WINDOW, gc.C_GREEN, gc.C_BLACK)
	gc.InitPair(UNFOCUSED_WINDOW, gc.C_WHITE, gc.C_BLACK)
	gc.InitPair(ERROR_WINDOW, gc.C_RED, gc.C_BLACK)

	return nil
}

func Loop() {
	for {
		rssReadersWin.Refresh()
		newsWin.Refresh()

		gc.Update()
		switch ch := rssReadersWin.GetChar(); ch {
		case 'q':
			rssReadersWin.Clear()
			newsWin.Clear()
			return
		case gc.KEY_TAB:
			activeWindow += 1
			if activeWindow > 1 {
				activeWindow = 0
			}
			panels[activeWindow].Top()
			for i, w := range windows {
				if i == activeWindow {
					w.Focus(FOCUSED_WINDOW)
				} else {
					w.Unfocus(UNFOCUSED_WINDOW)
				}
			}
		default:
			if activeWindow == RSS_READERS_WINDOW {
				rssReadersMenu.Driver(gc.DriverActions[ch])
			} else {
				newsMenu.Driver(gc.DriverActions[ch])
			}
		}
	}
}

func Main() {
	var rssReaders []db.RssReader
	var events []news.Event

	// Init DB
	if tdb, err = db.InitDB("./terminews.db"); err != nil {
		handleDBFatalError(err)
	}

	// Init UI components
	if err = InitUI(); err != nil {
		handleUIFatalError(err)
	}

	// Load RSS Readers data
	if rssReaders, err = tdb.GetRssReaders(); err != nil {
		handleDBFatalError(err)
	}
	if err = CreaterssReadersWindow(rssReaders); err != nil {
		handleUIFatalError(err)
	}

	// Load News
	if events, err = rss.Retrieve(rssReaders[0].Url); err != nil {
		handleNewsLoadError(rssReaders[0].Name)
	}
	if err = CreateNewsData(events, rssReaders[0].Name); err != nil {
		handleUIFatalError(err)
	}

	defer Free()

	Loop()
}
