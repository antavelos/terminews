package ui

import (
	"fmt"
	c "github.com/jroimartin/gocui"
)

type Displayer interface {
	Display() string
}

type List struct {
	*c.View
	name string
	Data Displayer
}

func CreateList(g *c.Gui, name string, x0, y0, x1, y1 int) (*List, error) {
	v, err := g.SetView(name, x0, y0, x1, y1)
	if err != nil && err != c.ErrUnknownView {
		return nil, err
	}

	list := &List{}
	list.View = v
	list.name = name
	list.SelBgColor = c.ColorWhite
	list.SelFgColor = c.ColorBlack
	list.BgColor = c.ColorBlack
	list.Autoscroll = true

	return list, nil
}

func (l *List) Focus(g *c.Gui) error {
	if _, err := g.SetCurrentView(l.name); err != nil {
		return err
	}
	return nil
}

func (l *List) SetItems(items []Displayer) error {
	l.Clear()
	for i, item := range items {
		if _, err := fmt.Fprintf(l.View, "%2d. %v\n", i+1, item.Display()); err != nil {
			return err
		}
	}

	return nil
}
