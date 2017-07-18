package ctrl

import (
	// "github.com/antavelos/terminews/db"s
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
)

var (
	Stdscr                             *gc.Window
	RssReadersWin, NewsWin             *ui.TWindow
	RssReadersMenu, NewsMenu           *ui.TMenu
	RssReadersMenuItems, NewsMenuItems []*ui.TMenuItem
	Windows                            []*ui.TWindow
	Panels                             []*gc.Panel
	ActiveWindow                       int
	err                                error
	maxX                               int
	maxY                               int
)

func CreateRssReadersWindow(rssReadersData []string) error {

	// MenuItems
	for i, d := range rssReadersData {
		item := &ui.TMenuItem{}
		if err = item.Create(i+1, d, d); err != nil {
			return err
		}
		RssReadersMenuItems = append(RssReadersMenuItems, item)
	}
	// Menu
	RssReadersMenu = &ui.TMenu{}
	if err = RssReadersMenu.Create(RssReadersMenuItems); err != nil {
		return err
	}
	// Window
	RssReadersWin = &ui.TWindow{}
	if err = RssReadersWin.Create("My RSS Readers", maxY, maxX/4, 0, 0); err != nil {
		return err
	}
	Windows = append(Windows, RssReadersWin)

	// attach RssReadersMenu on RssReadersWindow
	RssReadersWin.AttachMenu(RssReadersMenu)

	Panels = append(Panels, gc.NewPanel(RssReadersWin.Window))

	RssReadersMenu.Post()
	RssReadersWin.Refresh()

	return nil
}

func CreateNewsData(newsData []string) error {
	// MenuItems
	for i, d := range newsData {
		item := &ui.TMenuItem{}
		if err = item.Create(i+1, d, d); err != nil {
			return err
		}
		NewsMenuItems = append(NewsMenuItems, item)
	}
	// Menu
	NewsMenu = &ui.TMenu{}
	if err = NewsMenu.Create(NewsMenuItems); err != nil {
		return err
	}

	NewsWin = &ui.TWindow{}
	if err = NewsWin.Create("News from", maxY, (maxX*3)/4, 0, maxX/4); err != nil {
		return err
	}
	Windows = append(Windows, NewsWin)

	// attach NewsMenu on NewsWindow
	NewsWin.AttachMenu(NewsMenu)

	Panels = append(Panels, gc.NewPanel(NewsWin.Window))

	NewsMenu.Post()
	NewsWin.Refresh()

	return nil
}

// the order mstters!!!
func Free() {
	var tmi *ui.TMenuItem

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

func handleUIError(err error) {
	gc.End()
	log.Fatal(err)
}

func InitUI() error {

	Stdscr, err = gc.Init()
	if err != nil {
		return err
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

	return nil
}

func Main() {

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

	if err := InitUI(); err != nil {
		handleUIError(err)
	}

	if err := CreateRssReadersWindow(rssReadersData); err != nil {
		handleUIError(err)
	}

	if err := CreateNewsData(newsData); err != nil {
		handleUIError(err)
	}

	defer Free()

	RssReadersWin.Focus(FOCUSED_WINDOW)
	for {
		RssReadersWin.Refresh()
		NewsWin.Refresh()

		gc.Update()
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
					w.Focus(FOCUSED_WINDOW)
				} else {
					w.Unfocus(UNFOCUSED_WINDOW)
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
