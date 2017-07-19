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
	activeRssReadersMenuItem           int
	activeNewsMenuItem                 int
	err                                error
	maxX                               int
	maxY                               int
)

func CreateRssReadersMenu() error {
	// Menu
	rssReadersMenu = &ui.TMenu{}
	if err = rssReadersMenu.Create(); err != nil {
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

func CreateNewsMenu() error {
	// Menu
	newsMenu = &ui.TMenu{}
	if err = newsMenu.Create(); err != nil {
		return err
	}

	newsWin = &ui.TWindow{}
	if err = newsWin.Create("", maxY, (maxX*3)/4, 0, maxX/4); err != nil {
		return err
	}
	newsWin.SetHLine(maxY / 2)
	windows = append(windows, newsWin)

	// attach newsMenu on newsWindow
	newsWin.AttachMenu(newsMenu)

	panels = append(panels, gc.NewPanel(newsWin.Window))

	newsMenu.Post()
	newsWin.Refresh()

	return nil
}

func UpdateRssReadersMenuItems(rssReaders []db.RssReader) error {
	var displayers []ui.Displayer = make([]ui.Displayer, len(rssReaders))
	for i, rr := range rssReaders {
		displayers[i] = rr
	}
	rssReadersMenu.UnPost()
	tmis, err := rssReadersMenu.RefreshItems(displayers)
	if err != nil {
		return err
	}
	rssReadersMenu.Post()
	rssReadersMenuItems = tmis

	return nil
}

func UpdateNewsMenuItems(events []news.Event) error {
	var displayers []ui.Displayer = make([]ui.Displayer, len(events))
	for i, e := range events {
		displayers[i] = e
	}
	newsMenu.UnPost()
	tmis, err := newsMenu.RefreshItems(displayers)
	if err != nil {
		return err
	}
	newsMenuItems = tmis
	newsMenu.Post()

	return nil
}

func LoadNews(rssReader db.RssReader) {
	events, err := rss.Retrieve(rssReader.Url)
	if err != nil {
		handleUIFatalError(err)

		// handleNewsLoadError(err.Error())
	}
	if err = UpdateNewsMenuItems(events); err != nil {
		handleUIFatalError(err)
	}
	LoadNewsContent(events[0])
	newsWin.SetTitle(fmt.Sprintf("News from %v", rssReader.Name))
}

func LoadNewsContent(event news.Event) {
	authorLine := fmt.Sprintf("By %v", string(event.Author))
	publishedLine := fmt.Sprintf("Published on: %v", event.Published)
	linkLine := fmt.Sprintf("Link: %v", event.Link)
	summaryLine := fmt.Sprintf("%v", string(event.Description))
	halfway := (newsWin.H / 2)

	newsWin.SetHLine(halfway)
	newsWin.SetLine(authorLine, (halfway + 1), 2)
	newsWin.SetLine(publishedLine, (halfway + 2), 2)
	newsWin.SetLine(linkLine, (halfway + 3), 2)
	newsWin.SetLine(summaryLine, (halfway + 5), 2)
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
	newsWin.SetTitle(err.Error())
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

func onRssReadersWin() bool {
	return activeWindow == RSS_READERS_WINDOW
}

func onNewsWin() bool {
	return activeWindow == NEWS_WINDOW
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
			key := gc.Key(ch)
			if onRssReadersWin() {
				rssReadersMenu.Driver(gc.DriverActions[key])
			} else {
				newsMenu.Driver(gc.DriverActions[key])
			}
			if key == gc.KEY_UP {
				var active *int
				if onRssReadersWin() {
					active = &activeRssReadersMenuItem
				} else {
					active = &activeNewsMenuItem
				}
				if *active > 0 {
					*active -= 1
				}
				if onNewsWin() {
					item := newsMenuItems[activeNewsMenuItem]
					event := item.Data.(news.Event)
					LoadNewsContent(event)
				}
			}
			if key == gc.KEY_DOWN {
				var active *int
				var menuItems []*ui.TMenuItem
				if onRssReadersWin() {
					active = &activeRssReadersMenuItem
					menuItems = rssReadersMenuItems
				} else {
					active = &activeNewsMenuItem
					menuItems = newsMenuItems

				}
				if *active < len(menuItems)-1 {
					*active += 1
				}
				if onNewsWin() {
					item := newsMenuItems[activeNewsMenuItem]
					event := item.Data.(news.Event)
					LoadNewsContent(event)
				}
			}
			if key == gc.KEY_RETURN || key == gc.KEY_ENTER {
				if onRssReadersWin() {
					item := rssReadersMenuItems[activeRssReadersMenuItem]
					rssReader := item.Data.(db.RssReader)
					LoadNews(rssReader)
				}
			}
		}
	}
}

func Main() {
	var rssReaders []db.RssReader

	// Init DB
	if tdb, err = db.InitDB("./terminews.db"); err != nil {
		handleDBFatalError(err)
	}

	// Init UI components
	if err = InitUI(); err != nil {
		handleUIFatalError(err)
	}

	if err = CreateRssReadersMenu(); err != nil {
		handleUIFatalError(err)
	}

	if err = CreateNewsMenu(); err != nil {
		handleUIFatalError(err)
	}

	// Load RSS Readers data
	if rssReaders, err = tdb.GetRssReaders(); err != nil {
		handleDBFatalError(err)
	}
	if err = UpdateRssReadersMenuItems(rssReaders); err != nil {
		handleUIFatalError(err)
	}

	// Load News
	LoadNews(rssReaders[0])

	defer Free()

	Loop()
}
