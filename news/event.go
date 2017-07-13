package news

import (
	"fmt"
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
	return fmt.Sprintf("%v, %v", e.Title, e.Date)
}

type Events []Event

func (es Events) String() string {
	s := []string{}
	for i, e := range es {
		s = append(s, fmt.Sprintf("%v - %v", i+1, e.String()))
	}
	return strings.Join(s, "\n\n")
}
