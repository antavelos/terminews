package news

import (
	"fmt"
	"github.com/fatih/color"
	"strings"
	// "time"
)

type Event struct {
	Title       string
	Link        string
	Description string
	Date        string
}

func (e Event) String() string {
	yellow := color.New(color.FgYellow, color.Bold).SprintFunc()
	magenta := color.New(color.FgMagenta).SprintFunc()
	faint := color.New(color.Faint).SprintFunc()

	return fmt.Sprintf("%v, %v\n     %v\n     %v",
		yellow(e.Title), magenta(e.Date), faint(e.Link), e.Description)
}

func (e *Event) Display() string {
	return e.Title
}

type Events []Event

func (es Events) String() string {
	bold := color.New(color.Bold).SprintFunc()

	s := []string{}
	for i, e := range es {
		s = append(s, fmt.Sprintf("%v - %v", bold(i+1), e.String()))
	}
	return strings.Join(s, "\n\n")
}
