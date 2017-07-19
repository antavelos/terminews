package news

import (
	"fmt"
	"github.com/fatih/color"
)

type Event struct {
	Title       []rune
	Author      []rune
	Link        string
	Description []rune
	Published   string
}

func (e Event) String() string {
	yellow := color.New(color.FgYellow, color.Bold).SprintFunc()
	magenta := color.New(color.FgMagenta).SprintFunc()
	faint := color.New(color.Faint).SprintFunc()

	return fmt.Sprintf("%v, %v\n     %v\n     %v",
		yellow(string(e.Title)), magenta(e.Published), faint(e.Link), string(e.Description))
}

func (e Event) Display() string {
	return string(e.Title)
}
