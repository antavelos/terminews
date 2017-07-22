package ui

import (
	"fmt"
	c "github.com/jroimartin/gocui"
	_ "log"
)

type Displayer interface {
	Display() string
}

type List struct {
	*c.View
	name, title string
	Items       []Displayer
	currPage    int
}

func CreateList(g *c.Gui, name string, x0, y0, x1, y1 int) (*List, error) {
	v, err := g.SetView(name, x0, y0, x1, y1)
	if err != nil && err != c.ErrUnknownView {
		return nil, err
	}

	list := &List{}
	list.View = v
	list.name = name
	list.currPage = 1
	list.SelBgColor = c.ColorWhite
	list.SelFgColor = c.ColorBlack
	list.Autoscroll = true

	return list, nil
}

func (l *List) height() int {
	_, y := l.Size()
	return y - 1
}

func (l *List) width() int {
	x, _ := l.Size()
	return x
}

func (l *List) length() int {
	return len(l.Items)
}

func (l *List) pages() int {
	d := l.length() / l.height()

	if l.length()%l.height() > 0 {
		d++
	}

	return d
}

func (l *List) nextPage() (int, int) {
	if l.currPage < l.pages() {
		l.currPage++
	}
	return l.getPage(l.currPage)
}

func (l *List) prevPage() (int, int) {
	if l.currPage > 1 {
		l.currPage--
	}
	return l.getPage(l.currPage)
}

func (l *List) firstPage() (int, int) {
	l.currPage = 1
	return l.getPage(l.currPage)
}

func (l *List) lastPage() (int, int) {
	l.currPage = l.pages()
	return l.getPage(l.currPage)
}

func (l *List) hasNextPage() bool {
	return (l.pages() - l.currPage) > 0
}

func (l *List) hasPrevPage() bool {
	return l.currPage > 1
}

func (l *List) getPage(p int) (int, int) {
	var start, end int
	start = (p - 1) * l.height()
	if p == l.pages() {
		end = l.length()
	} else {
		end = start + l.height()
	}
	return start, end
}

func (l *List) Focus(g *c.Gui) error {
	if _, err := g.SetCurrentView(l.name); err != nil {
		return err
	}
	l.Highlight = true
	return nil
}

func (l *List) Unfocus() {
	l.Highlight = false
}

func (l *List) SetTitle(title string) {
	l.title = title
	l.Title = fmt.Sprintf(" %d/%d - %v ", l.currPage, l.pages(), title)
}

func (l *List) SetItems(data []Displayer) error {
	l.Items = data

	var start, end int
	if l.length() > l.height() {
		start, end = 0, l.height()
	} else {
		start, end = 0, l.length()
	}
	l.display(start, end)
	l.currPage = 1

	return nil
}

func (l *List) displayItem(i int) string {
	item := l.Items[i]
	return fmt.Sprintf("%2d. %v", i+1, item.Display())
}

func (l *List) display(start, end int) error {
	l.Clear()
	for i := start; i < end; i++ {
		if _, err := fmt.Fprintf(l.View, "%v\n", l.displayItem(i)); err != nil {
			return err
		}
	}
	l.SetTitle(l.title)
	return nil
}

func (l *List) MoveDown() error {
	x, y := l.Cursor()

	if y == l.height()-1 {
		y = 0
		if l.hasNextPage() {
			l.display(l.nextPage())
		} else {
			l.display(l.firstPage())
		}
	} else {
		y++
	}
	if err := l.SetCursor(x, y); err != nil {
		return err
	}

	return nil
}

func (l *List) MoveUp() error {
	x, y := l.Cursor()
	if y == 0 {
		y = l.height() - 1
		if l.hasPrevPage() {
			l.display(l.prevPage())
		} else {
			l.display(l.lastPage())
		}
	} else {
		y--
	}

	if err := l.SetCursor(x, y); err != nil {
		return err
	}

	return nil
}

func (l *List) CurrentItem() Displayer {
	_, y := l.Cursor()

	start, end := l.getPage(l.currPage)
	data := l.Items[start:end]

	return data[y]
}

func (l *List) CurrentPage() int {
	return l.currPage
}
func (l *List) ResetCursor() {
	l.SetCursor(0, 0)
}
