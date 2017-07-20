package ctrl

import (
	"bytes"
	"fmt"
	"log"

	// Both TUI packages are abbreviated to avoid making the code
	// overly verbose.
	"github.com/fatih/color"
	c "github.com/jroimartin/gocui"
	"github.com/pkg/errors"
)

// Items to fill the list with.
var listItems = []string{
	"Line 1",
	"Line 2",
	"Line 3",
	"Line 4",
	"Line 5",
}

// var tdb *db.TDB
var err error

// func handleDBFatalError(err error) {
// 	tdb.Close()
// 	log.Fatal(err)
// }

func spaces(n int) string {
	var s bytes.Buffer
	for i := 0; i < n; i++ {
		s.WriteString(" ")
	}
	return s.String()

}

// Set up the widgets and run the event loop.
func Main() {
	// Create a new GUI.
	g, err := c.NewGui(c.OutputNormal)
	if err != nil {
		log.Println("Failed to create a GUI:", err)
		return
	}
	defer g.Close()

	// g.Cursor = true

	g.SetManagerFunc(layout)

	err = g.SetKeybinding("", c.KeyCtrlC, c.ModNone, quit)
	if err != nil {
		log.Println("Could not set key binding:", err)
		return
	}

	tw, th := g.Size()

	lw := (tw * 3) / 10
	oh := (th * 6) / 10

	lv, err := g.SetView("rssreaders", 0, 0, lw, th-1)
	if err != nil && err != c.ErrUnknownView {
		log.Println("Failed to create rssreaders view:", err)
		return
	}
	wb := color.New(color.FgBlack, color.BgWhite)

	lv.Title = " RSS Readers "
	lv.FgColor = c.ColorCyan

	// Then the output view.
	ov, err := g.SetView("news", lw+1, 0, tw-1, oh)
	if err != nil && err != c.ErrUnknownView {
		log.Println("Failed to create news view:", err)
		return
	}
	ov.Title = "News from ..."
	ov.FgColor = c.ColorGreen
	// Let the view scroll if the output exceeds the visible area.
	ov.Autoscroll = true
	ov.Wrap = true

	// And finally the input view.
	iv, err := g.SetView("summary", lw+1, oh+1, tw-1, th-1)
	if err != nil && err != c.ErrUnknownView {
		log.Println("Failed to create summary view:", err)
		return
	}
	iv.Title = "Summary"
	iv.FgColor = c.ColorYellow
	// The input view shall be editable.

	// Make the enter key copy the input to the output.
	err = g.SetKeybinding("input", c.KeyEnter, c.ModNone, func(g *c.Gui, iv *c.View) error {
		// We want to read the view's buffer from the beginning.
		iv.Rewind()

		// Get the output view via its name.
		ov, e := g.View("output")
		if e != nil {
			log.Println("Cannot get output view:", e)
			return e
		}
		// Thanks to views being an io.Writer, we can simply Fprint to a view.
		_, e = fmt.Fprint(ov, iv.Buffer())
		if e != nil {
			log.Println("Cannot print to output view:", e)
		}
		// Clear the input view
		iv.Clear()
		// Put the cursor back to the start.
		e = iv.SetCursor(0, 0)
		if e != nil {
			log.Println("Failed to set cursor:", e)
		}
		return e

	})
	if err != nil {
		log.Println("Cannot bind the enter key:", err)
	}

	// Fill the list view.
	for _, s := range listItems {
		// Again, we can simply Fprint to a view.
		_, err = wb.Fprintln(lv, s)
		if err != nil {
			log.Println("Error writing to the list view:", err)
			return
		}
	}
	// Set the focus to the input view.
	_, err = g.SetCurrentView("rssreaders")
	if err != nil {
		log.Println("Cannot set focus to input view:", err)
	}

	// Start the main loop.
	err = g.MainLoop()
	log.Println("Main loop has finished:", err)
}

// The layout handler calculates all sizes depending
// on the current terminal size.
func layout(g *c.Gui) error {
	// Get the current terminal size.
	tw, th := g.Size()

	lw := (tw * 3) / 10
	oh := (th * 6) / 10

	_, err := g.SetView("rssreaders", 0, 0, lw, th-1)
	if err != nil {
		return errors.Wrap(err, "Cannot update list view")
	}
	_, err = g.SetView("news", lw+1, 0, tw-1, oh)
	if err != nil {
		return errors.Wrap(err, "Cannot update output view")
	}
	_, err = g.SetView("summary", lw+1, oh+1, tw-1, th-1)
	if err != nil {
		return errors.Wrap(err, "Cannot update input view.")
	}
	return nil
}

// `quit` is a handler that gets bound to Ctrl-C.
// It signals the main loop to exit.
func quit(g *c.Gui, v *c.View) error {
	return c.ErrQuit
}
