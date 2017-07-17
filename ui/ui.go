package ui

import (
	"fmt"
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
)

type TMenuItem struct {
	*gc.MenuItem
	id   int
	name string
	desc string
}

func (tmi *TMenuItem) Create(id int, name string, desc string) error {
	tmi.id = id
	tmi.name = name
	tmi.desc = desc
	value := fmt.Sprintf("%3d. %v", id, name)
	item, err := gc.NewItem(value, "")
	if err != nil {
		return err
	}
	tmi.MenuItem = item

	return nil
}

type TMenu struct {
	*gc.Menu
	items []*TMenuItem
}

func (tm *TMenu) Create(tmis []*TMenuItem) error {
	tm.items = tmis
	gcMenuItems := make([]*gc.MenuItem, len(tmis))
	for i, tmi := range tmis {
		gcMenuItems[i] = tmi.MenuItem
	}
	menu, err := gc.NewMenu(gcMenuItems)
	if err != nil {
		return err
	}
	tm.Menu = menu

	return nil
}

type TWindow struct {
	*gc.Window
	title      string
	h, w, y, x int
}

func (tw *TWindow) Create(title string, h, w, y, x int) error {
	win, err := gc.NewWindow(h, w, y, x)
	if err != nil {
		return err
	}
	tw.Window = win
	tw.h = h
	tw.w = w
	tw.y = y
	tw.x = x

	tw.Keypad(true)
	tw.SetContour()
	tw.SetTitle(title)

	return nil
}

func (tw *TWindow) SetTitle(title string) {
	tw.title = title
	// _, mx := tw.MaxYX()
	tw.MovePrint(1, (tw.w/2)-(len(title)/2), title)
}

func (tw *TWindow) SetContour() {
	tw.Box(0, 0)
	tw.MoveAddChar(2, 0, gc.ACS_LTEE)
	tw.HLine(2, 1, gc.ACS_HLINE, tw.w-2)
	tw.MoveAddChar(2, tw.w-1, gc.ACS_RTEE)
}

func (tw *TWindow) Focus() {
	tw.ColorOn(FOCUSED_WINDOW)
	tw.SetContour()
	tw.ColorOff(FOCUSED_WINDOW)
}

func (tw *TWindow) Unfocus() {
	tw.ColorOn(UNFOCUSED_WINDOW)
	tw.SetContour()
	tw.ColorOff(UNFOCUSED_WINDOW)
}

func (tw *TWindow) attachMenu(tm *TMenu) {
	tm.Menu.SetWindow(tw.Window)
	tm.Menu.SubWindow(tw.Derived(tw.h-6, tw.w-4, 4, 2))
	tm.Menu.Format(tw.h-6, 1)
	tm.Menu.Mark("")
}

var (
	Stdscr                             *gc.Window
	RssReadersWin, NewsWin             *TWindow
	RssReadersMenu, NewsMenu           *TMenu
	RssReadersMenuItems, NewsMenuItems []*TMenuItem
	Windows                            []*TWindow
	Panels                             []*gc.Panel
	ActiveWindow                       int
	err                                error
	maxX                               int
	maxY                               int
)

// the order mstters!!!
func Free() {
	var tmi *TMenuItem

	RssReadersMenu.UnPost()
	NewsMenu.UnPost()

	for _, tmi = range RssReadersMenuItems {
		tmi.MenuItem.Free()
	}
	for _, tmi = range NewsMenuItems {
		tmi.MenuItem.Free()
	}

	RssReadersMenu.Free()
	NewsMenu.Free()

	gc.End()
}

func Init() {

	Stdscr, err = gc.Init()
	if err != nil {
		handleError(err)
	}

	maxY, maxX = Stdscr.MaxYX()
	// Initial configuration
	gc.StartColor()
	gc.Raw(true)
	gc.Echo(false)
	gc.Cursor(0)
	Stdscr.Keypad(true)
	Stdscr.Clear()

	// Define color combinations
	gc.InitPair(FOCUSED_WINDOW, gc.C_GREEN, gc.C_BLACK)
	gc.InitPair(UNFOCUSED_WINDOW, gc.C_WHITE, gc.C_BLACK)
}

func CreateRssReadersWindow(rssReadersData []string) {

	// MenuItems
	for i, d := range rssReadersData {
		item := &TMenuItem{}
		if err = item.Create(i+1, d, d); err != nil {
			handleError(err)
		}
		RssReadersMenuItems = append(RssReadersMenuItems, item)
	}
	// Menu
	RssReadersMenu = &TMenu{}
	if err = RssReadersMenu.Create(RssReadersMenuItems); err != nil {
		handleError(err)
	}
	// Window
	RssReadersWin = &TWindow{}
	if err = RssReadersWin.Create("My RSS Readers", maxY, maxX/4, 0, 0); err != nil {
		handleError(err)
	}
	Windows = append(Windows, RssReadersWin)

	// attach RssReadersMenu on RssReadersWindow
	RssReadersWin.attachMenu(RssReadersMenu)

	Panels = append(Panels, gc.NewPanel(RssReadersWin.Window))

	RssReadersMenu.Post()
	RssReadersWin.Refresh()
}

func CreateNewsData(newsData []string) {
	// MenuItems
	for i, d := range newsData {
		item := &TMenuItem{}
		if err = item.Create(i+1, d, d); err != nil {
			handleError(err)
		}
		NewsMenuItems = append(NewsMenuItems, item)
	}
	// Menu
	NewsMenu = &TMenu{}
	if err = NewsMenu.Create(NewsMenuItems); err != nil {
		handleError(err)
	}

	NewsWin = &TWindow{}
	if err = NewsWin.Create("News from", maxY, (maxX*3)/4, 0, maxX/4); err != nil {
		handleError(err)
	}
	Windows = append(Windows, NewsWin)

	// attach NewsMenu on NewsWindow
	NewsWin.attachMenu(NewsMenu)

	Panels = append(Panels, gc.NewPanel(NewsWin.Window))

	NewsMenu.Post()
	NewsWin.Refresh()
}

func handleError(err error) {
	gc.End()
	log.Fatal(err)
}

func Ui() {

	newsData := []string{
		"Life inside ISIS bride camp: Fighting, sex obsessed fighters, 'jihadi Tinder'",
		"Guns and money: Why US' top North Korea diplomat is in Southeast Asia",
		"10 days of horror: Grisly killings stun a small town",
		"Secret Service refutes Trump lawyer remarks",
		"UAE denies report it orchestrated Qatar hack",
		"Russia rejects any US conditions for return of seized compounds",
		"Jordanian soldier gets life for killing US troops",
		"Merkel rules out limiting number of refugees in Germany",
		"Australian woman in US fatally shot by police officer",
		"65 arrested in Europe-wide horsemeat scam",
		"Gel is 5 times stronger than steel",
		"Family finds clues to teen's suicide in blue whale paintings",
		"Ex-Mexican president banned from Venezuela",
		"Where North Korea's elite go for banned luxury goods",
		"BBC makes history with 'Doctor Who' casting",
		"Duchess of Cornwall speaks to CNN in rare interview",
		"Actor Martin Landau dies at 89",
		"Columbia University settles with student accused of sexual assault",
		"Delta hits back after tweetstorm mix-up",
		"What China's new GDP numbers reveal",
		"Horror genre godfather George Romero dead at 77",
		"Flash flood kills 7 from family",
		"16 pilgrims die when bus falls into gorge",
		"Iran sentences American to 10 years on spying conviction ",
		"'No one ever helps us': Life after escaping conflict in Myanmar",
		"$1 million in pot found in brand-new cars",
		"Can artificial sweeteners cause weight gain?",
		"F1 owners: British Grand Prix will stay",
		"Mediterranean style diet may prevent dementia",
		"Greatest golf links in the world",
		"What did we learn from Mayweather vs. McGregor traveling circus?",
		"That huge iceberg should freak you out. Here's why",
		"Ubud: Inside Bali's cultural epicenter",
		"Why Guinness tastes different in Africa",
		"Roger Federer beats Cilic to clinch historic eighth Wimbledon title",
		"Pink Floyd star defends his anti-Trump tour",
		"Hamilton wins record 5th British GP",
		"Who will win flying car race?",
		"A new breed of supercar",
		"Every airport should be like this",
		"Republicans delay health vote after McCain has surgery",
		"'Walking Dead' stuntman dies after fall on set",
		"Why is ocean being pumped into desert?",
		"Repairman gets stuck in ATM",
		"Airline sends rapper's dog to wrong city",
		"How to design the car of your dreams in VR ",
		"How ISIS changed Iraqi schools",
		"Is the fall of Mosul the fall of ISIS?",
		"Why US' top N. Korea diplomat is in  Asia",
	}
	// Create RSS Readers menu
	rssReadersData := []string{
		"CNN World",
		"BBC",
		"NBC",
		"Reuters",
	}

	Init()

	CreateRssReadersWindow(rssReadersData)
	CreateNewsData(newsData)

	defer Free()

	RssReadersWin.Focus()
	for {
		RssReadersWin.Refresh()
		NewsWin.Refresh()

		gc.Update()
		// gc.UpdatePanels()
		switch ch := RssReadersWin.GetChar(); ch {
		case 'q':
			RssReadersWin.Clear()
			NewsWin.Clear()
			return
		case gc.KEY_TAB:
			ActiveWindow += 1
			if ActiveWindow > 1 {
				ActiveWindow = 0
			}
			Panels[ActiveWindow].Top()
			for i, w := range Windows {
				if i == ActiveWindow {
					w.Focus()
				} else {
					w.Unfocus()
				}
			}
		default:
			if ActiveWindow == RSS_READERS_WINDOW {
				RssReadersMenu.Driver(gc.DriverActions[ch])
			} else {
				NewsMenu.Driver(gc.DriverActions[ch])
			}
		}
	}
}
