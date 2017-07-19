package ui

import (
	"fmt"
	gc "github.com/rthornton128/goncurses"
)

type Displayer interface {
	Display() string
}

type TMenuItem struct {
	*gc.MenuItem
	inc  int
	Data Displayer
}

func (tmi *TMenuItem) Create(inc int, data Displayer) error {
	value := fmt.Sprintf(`%3d. %+q`, inc, data.Display())
	item, err := gc.NewItem(value, "")
	if err != nil {
		return err
	}
	tmi.inc = inc
	tmi.Data = data
	tmi.MenuItem = item

	return nil
}

type TMenu struct {
	*gc.Menu
	items []*TMenuItem
}

func (tm *TMenu) Create() error {
	menu, err := gc.NewMenu([]*gc.MenuItem{})
	if err != nil {
		return err
	}
	tm.Menu = menu

	return nil
}

func (tm *TMenu) RefreshItems(dataItems []Displayer) ([]*TMenuItem, error) {
	var tmis []*TMenuItem
	var gcItems []*gc.MenuItem

	for _, gcmi := range tm.Items() {
		gcmi.Free()
	}

	for i, data := range dataItems {
		tmi := &TMenuItem{}
		if err := tmi.Create(i+1, data); err != nil {
			return nil, err
		}
		tmis = append(tmis, tmi)
	}
	for _, tmi := range tmis {
		gcItems = append(gcItems, tmi.MenuItem)
	}
	if err := tm.SetItems(gcItems); err != nil {
		return nil, err
	}
	tm.items = tmis

	return tmis, nil
}

type TWindow struct {
	*gc.Window
	title      string
	H, W, Y, X int
}

func (tw *TWindow) Create(title string, h, w, y, x int) error {
	win, err := gc.NewWindow(h, w, y, x)
	if err != nil {
		return err
	}
	tw.Window = win
	tw.H = h
	tw.W = w
	tw.Y = y
	tw.X = x

	tw.Keypad(true)
	tw.Box(0, 0)
	tw.SetHLine(2)
	tw.SetTitle(title)

	return nil
}

func (tw *TWindow) SetTitle(title string) {
	tw.title = title
	tw.SetLine(title, 1, (tw.W/2)-(len(title)/2))
}

func (tw *TWindow) SetHLine(y int) {
	tw.HLine(y, 1, gc.ACS_HLINE, tw.W-2)
}

func (tw *TWindow) SetLine(s string, y, x int) {
	tw.MovePrint(y, x, s)
}

func (tw *TWindow) Focus(cc int16) {
	tw.ColorOn(cc)
	tw.Box(0, 0)
	tw.ColorOff(cc)
}

func (tw *TWindow) Unfocus(cc int16) {
	tw.ColorOn(cc)
	tw.Box(0, 0)
	tw.ColorOff(cc)
}

func (tw *TWindow) AttachMenu(tm *TMenu) {
	tm.Menu.SetWindow(tw.Window)
	tm.Menu.SubWindow(tw.Derived(tw.H-4, tw.W-4, 4, 2))
	tm.Menu.Format(tw.H-4, 1)
	tm.Menu.Mark("")
}
