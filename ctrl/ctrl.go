package ctrl

import (
	// "github.com/antavelos/terminews/db"s
	"github.com/antavelos/terminews/ui"
	gc "github.com/rthornton128/goncurses"
)

func Main() {
	ui.Init()

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

	ui.CreateRssReadersWindow(rssReadersData)
	ui.CreateNewsData(newsData)

	defer ui.Free()

	ui.RssReadersWin.Focus()
	for {
		ui.RssReadersWin.Refresh()
		ui.NewsWin.Refresh()

		gc.Update()
		switch ch := ui.RssReadersWin.GetChar(); ch {
		case 'q':
			ui.RssReadersWin.Clear()
			ui.NewsWin.Clear()
			return
		case gc.KEY_TAB:
			ui.ActiveWindow += 1
			if ui.ActiveWindow > 1 {
				ui.ActiveWindow = 0
			}
			ui.Panels[ui.ActiveWindow].Top()
			for i, w := range ui.Windows {
				if i == ui.ActiveWindow {
					w.Focus()
				} else {
					w.Unfocus()
				}
			}
		default:
			if ui.ActiveWindow == ui.RSS_READERS_WINDOW {
				ui.RssReadersMenu.Driver(gc.DriverActions[ch])
			} else {
				ui.NewsMenu.Driver(gc.DriverActions[ch])
			}
		}
	}
}
