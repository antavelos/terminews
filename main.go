/*
   Terminews is a terminal based (TUI) RSS feed manager.
   Copyright (C) 2017  Alexandros Ntavelos, a[dot]ntavelos[at]gmail[dot]com

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/
package main

import (
	"log"
	"sync"

	"github.com/antavelos/terminews/db"
	"github.com/fatih/color"
	c "github.com/jroimartin/gocui"
)

const (
	SITES_VIEW   = "rssreaders"
	NEWS_VIEW    = "news"
	SUMMARY_VIEW = "summary"
	PROMPT_VIEW  = "prompt"
)

var (
	Lists map[string]*List
	tdb   *db.TDB
	// g        *c.Gui
	// err      error
	sitesList *List
	newsList  *List
	summary   *c.View
	curW      int
	curH      int
	Bold      *color.Color
	wg        sync.WaitGroup
)

func handleFatalError(msg string, err error) {
	// tdb.Close()
	// g.Close()
	log.Fatal(msg, err)
}

func relSize(g *c.Gui) (int, int) {
	tw, th := g.Size()

	return (tw * 3) / 10, (th * 70) / 100
}

// The layout handler calculates all sizes depending
// on the current terminal size.
func layout(g *c.Gui) error {
	// Get the current terminal size.
	tw, th := g.Size()

	// Get the relative size of the views
	rw, rh := relSize(g)

	_, err := g.SetView(SITES_VIEW, 0, 0, rw, th-1)
	if err != nil {
		handleFatalError("Cannot update sites view", err)
	}

	_, err = g.SetView(NEWS_VIEW, rw+1, 0, tw-1, rh)
	if err != nil {
		handleFatalError("Cannot update news view", err)
	}

	_, err = g.SetView(SUMMARY_VIEW, rw+1, rh+1, tw-1, th-1)
	if err != nil {
		handleFatalError("Cannot update summary view.", err)
	}
	updateSummary()

	if _, err = g.View(PROMPT_VIEW); err == nil {
		_, err = g.SetView(PROMPT_VIEW, tw/6, (th/2)-1, (tw*5)/6, (th/2)+1)
		if err != nil && err != c.ErrUnknownView {
			return err
		}
	}

	if curW != tw || curH != th {
		sitesList.Draw()
		newsList.Draw()
		curW = tw
		curH = th
	}

	return nil
}

func main() {

	var v *c.View
	var err error

	Bold = color.New(color.Bold)
	Lists = make(map[string]*List)

	// Init DB
	if tdb, err = db.InitDB("./term.db"); err != nil {
		tdb.Close()
		log.Fatal(err)
	}
	defer tdb.Close()

	// Create a new GUI.
	g, err := c.NewGui(c.OutputNormal)
	if err != nil {
		handleFatalError("Failed to initialize GUI", err)
	}
	defer g.Close()

	g.SelFgColor = c.ColorGreen | c.AttrBold
	g.BgColor = c.ColorDefault
	g.Highlight = true

	g.SetManagerFunc(layout)

	curW, curH = g.Size()
	rw, rh := relSize(g)

	// Sites List
	v, err = g.SetView(SITES_VIEW, 0, 0, rw, curH-1)
	if err != nil && err != c.ErrUnknownView {
		handleFatalError("Failed to create sites list:", err)

	}
	sitesList = CreateList(v)

	//
	v, err = g.SetView(NEWS_VIEW, rw+1, 0, curW-1, rh)
	if err != nil && err != c.ErrUnknownView {
		handleFatalError(" Failed to create news list:", err)
	}
	newsList = CreateList(v)
	newsList.SetTitle("No news yet...")

	// Summary view
	summary, err = g.SetView(SUMMARY_VIEW, rw+1, rh+1, curW-1, curH-1)
	if err != nil && err != c.ErrUnknownView {
		handleFatalError("Failed to create summary view:", err)
	}
	summary.Title = " Summary "
	summary.Wrap = true

	loadSites()
	sitesList.Focus(g)

	g.SetRune(curW/2, curH/2, rune('a'), c.ColorDefault|c.AttrBold, c.ColorDefault)

	addKeybinding(g, "", c.KeyCtrlN, c.ModNone, addSite)
	addKeybinding(g, "", c.KeyDelete, c.ModNone, deleteEntry)
	addKeybinding(g, NEWS_VIEW, c.KeyCtrlB, c.ModNone, addBookmark)
	addKeybinding(g, "", c.KeyCtrlC, c.ModNone, quit)
	addKeybinding(g, "", c.KeyCtrlB, c.ModAlt, loadBookmarks)
	addKeybinding(g, "", c.KeyTab, c.ModNone, switchView)
	addKeybinding(g, "", c.KeyArrowUp, c.ModNone, listUp)
	addKeybinding(g, "", c.KeyArrowDown, c.ModNone, listDown)
	addKeybinding(g, "", c.KeyPgup, c.ModNone, listPgUp)
	addKeybinding(g, "", c.KeyPgdn, c.ModNone, listPgDown)
	addKeybinding(g, "", c.KeyEnter, c.ModNone, onEnter)
	addKeybinding(g, PROMPT_VIEW, c.KeyCtrlQ, c.ModNone, removePrompt)
	addKeybinding(g, "", c.KeyCtrlF, c.ModNone, find)

	err = g.MainLoop()
	log.Println("terminews exited unexpectedly: ", err)

}
