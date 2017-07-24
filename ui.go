package main

import (
	"bytes"
	"fmt"

	c "github.com/jroimartin/gocui"
)

type Page struct {
	offset, limit int
}

type List struct {
	*c.View
	title       string
	items       []interface{}
	pages       []Page
	currPageIdx int
}

func CreateList(v *c.View) *List {
	list := &List{}
	list.View = v
	list.SelBgColor = c.ColorWhite
	list.SelFgColor = c.ColorBlack
	list.Autoscroll = true

	return list
}

func (l *List) IsEmpty() bool {
	return l.length() == 0
}

func (l *List) Focus(g *c.Gui) error {
	l.Highlight = true
	_, err := g.SetCurrentView(l.Name())

	return err
}

func (l *List) Unfocus() {
	l.Highlight = false
}

func (l *List) SetTitle(title string) {
	l.title = title
	l.Title = fmt.Sprintf(" %d/%d - %v ", l.currPageIdx+1, l.pagesNum(), title)
}

func (l *List) SetItems(data []interface{}) error {
	l.items = data
	return l.Draw()
}

func (l *List) Draw() error {
	if l.IsEmpty() {
		return nil
	}
	l.currPageIdx = 0
	l.pages = []Page{}
	for offset := 0; offset < l.length(); offset += l.height() {
		limit := l.height()
		if offset+limit > l.length() {
			limit = l.length() % l.height()
		}
		l.pages = append(l.pages, Page{offset, limit})
	}
	return l.displayPage(l.currPageIdx)
}

func (l *List) MoveDown() error {
	y := l.currentCursorY() + 1
	if l.atBottomOfPage() {
		y = 0
		if l.hasMultiplePages() {
			l.displayPage(l.nextPageIdx())
		}
	}
	return l.SetCursor(0, y)
}

func (l *List) MoveUp() error {
	y := l.currentCursorY() - 1
	if l.atTopOfPage() {
		y = l.pages[l.prevPageIdx()].limit - 1
		if l.hasMultiplePages() {
			l.displayPage(l.prevPageIdx())
		}
	}

	return l.SetCursor(0, y)
}

func (l *List) MovePgDown() error {
	l.displayPage(l.nextPageIdx())

	return l.SetCursor(0, 0)
}

func (l *List) MovePgUp() error {
	l.displayPage(l.prevPageIdx())

	return l.SetCursor(0, 0)
}

func (l *List) CurrentItem() interface{} {
	page := l.currPage()
	data := l.items[page.offset : page.offset+page.limit]

	return data[l.currentCursorY()]
}

func (l *List) ResetCursor() {
	l.SetCursor(0, 0)
}

func (l *List) currentCursorY() int {
	_, y := l.Cursor()

	return y
}

func (l *List) currPage() Page {
	return l.pages[l.currPageIdx]
}

func (l *List) height() int {
	_, y := l.Size()

	return y - 1
}

func (l *List) width() int {
	x, _ := l.Size()

	return x - 1
}

func (l *List) length() int {
	return len(l.items)
}

func (l *List) pagesNum() int {
	return len(l.pages)
}

func (l *List) nextPageIdx() int {
	return (l.currPageIdx + 1) % l.pagesNum()
}

func (l *List) prevPageIdx() int {
	pidx := (l.currPageIdx - 1) % l.pagesNum()
	if l.currPageIdx == 0 {
		pidx = l.pagesNum() - 1
	}
	return pidx
}

func (l *List) displayItem(i int) string {
	item := fmt.Sprint(l.items[i])
	sp := spaces(l.width() - len(item) - 3)
	return fmt.Sprintf("%2d. %v%v", i+1, item, sp)
}

func (l *List) displayPage(p int) error {
	l.currPageIdx = p
	l.Clear()
	l.SetTitle(l.title)
	page := l.currPage()
	for i := page.offset; i < page.offset+page.limit; i++ {
		if _, err := fmt.Fprintf(l.View, "%v\n", l.displayItem(i)); err != nil {
			return err
		}
	}

	return nil
}

func (l *List) atBottomOfPage() bool {
	return l.currentCursorY() == l.currPage().limit-1
}

func (l *List) atTopOfPage() bool {
	return l.currentCursorY() == 0
}

func (l *List) hasMultiplePages() bool {
	return l.pagesNum() > 1
}

func spaces(n int) string {
	var s bytes.Buffer
	for i := 0; i < n; i++ {
		s.WriteString(" ")
	}
	return s.String()
}
