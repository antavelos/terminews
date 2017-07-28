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
	"os"
	"os/user"
	"path"

	"github.com/antavelos/terminews/db"
	"github.com/fatih/color"
	c "github.com/jroimartin/gocui"
)

const (
	SITES_VIEW   = "rssreaders"
	NEWS_VIEW    = "news"
	SUMMARY_VIEW = "summary"
	PROMPT_VIEW  = "prompt"
	CONTENT_VIEW = "content"
)

var (
	Lists          map[string]*List
	tdb            *db.TDB
	sitesList      *List
	newsList       *List
	contentList    *List
	summary        *c.View
	CurrentContent []string
	curW           int
	curH           int
	Bold           *color.Color
)

// relSize calculates the  sizes of the sites view width
// and the news view height in relation to the current terminal size
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
		log.Fatal("Cannot update sites view", err)
	}

	_, err = g.SetView(NEWS_VIEW, rw+1, 0, tw-1, rh)
	if err != nil {
		log.Fatal("Cannot update news view", err)
	}

	_, err = g.SetView(SUMMARY_VIEW, rw+1, rh+1, tw-1, th-1)
	if err != nil {
		log.Fatal("Cannot update summary view.", err)
	}
	updateSummary()

	if _, err = g.View(PROMPT_VIEW); err == nil {
		_, err = g.SetView(PROMPT_VIEW, tw/6, (th/2)-1, (tw*5)/6, (th/2)+1)
		if err != nil && err != c.ErrUnknownView {
			return err
		}
	}

	if _, err = g.View(CONTENT_VIEW); err == nil {
		_, err = g.SetView(CONTENT_VIEW, tw/8, th/8, (tw*7)/8, (th*7)/8)
		if err != nil && err != c.ErrUnknownView {
			return err
		}
	}

	if curW != tw || curH != th {
		sitesList.ResetPages()
		sitesList.Draw()
		newsList.ResetPages()
		newsList.Draw()
		if contentList != nil {
			contentList.Reset()
			UpdateContent(g, CurrentContent)
			// contentList.Draw()
		}
		curW = tw
		curH = th
	}

	return nil
}

// getappDir creates if not exists the app directory where the sqlite db file
// as well as the log file will be stored. In case of failure the current dir
// will be used.
func getAppDir() (string, error) {
	usr, _ := user.Current()
	dir := path.Join(usr.HomeDir, ".terminews")
	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			if oserr := os.Mkdir(dir, 0666); oserr != nil {
				return "", oserr
			}
		} else {
			return "", err
		}
	}
	return dir, nil
}

func main() {

	var v *c.View
	var err error

	Bold = color.New(color.Bold)
	Lists = make(map[string]*List)

	appDir, err := getAppDir()
	if err != nil {
		panic("Could not set up app directory.")
	}

	// Setup logging
	logfile := path.Join(appDir, "terminews.log")
	f, err := os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open logfile", err)
	}
	defer f.Close()
	log.SetOutput(f)

	// Init DB
	if tdb, err = db.InitDB(appDir); err != nil {
		log.Fatal("Failed to initialize DB", err)
	}
	defer tdb.Close()

	// Create a new GUI.
	g, err := c.NewGui(c.OutputNormal)
	if err != nil {
		log.Fatal("Failed to initialize GUI", err)
	}
	defer g.Close()

	// some basic configuration
	g.SelFgColor = c.ColorGreen | c.AttrBold
	g.BgColor = c.ColorDefault
	g.Highlight = true

	// setup the layout
	g.SetManagerFunc(layout)

	// the current actual size of the terminal
	curW, curH = g.Size()

	// rw is the relative width of the sites view
	// rh is the relative height of the news view
	rw, rh := relSize(g)

	// Setup the initial layout
	// Sites List
	v, err = g.SetView(SITES_VIEW, 0, 0, rw, curH-1)
	if err != nil && err != c.ErrUnknownView {
		log.Fatal("Failed to create sites list:", err)
	}
	sitesList = CreateList(v, true)
	sitesList.Focus(g)

	// it loads the existing sites if any at the beginning
	g.Execute(func(g *c.Gui) error {
		if err := loadSites(); err != nil {
			log.Fatal("Error while loading sites", err)
		}
		log.Print("Loaded initial sites")
		return nil
	})

	// News list
	v, err = g.SetView(NEWS_VIEW, rw+1, 0, curW-1, rh)
	if err != nil && err != c.ErrUnknownView {
		log.Fatal(" Failed to create news list:", err)
	}
	newsList = CreateList(v, true)
	newsList.SetTitle("No news yet...")

	// Summary view
	summary, err = g.SetView(SUMMARY_VIEW, rw+1, rh+1, curW-1, curH-1)
	if err != nil && err != c.ErrUnknownView {
		log.Fatal("Failed to create summary view:", err)
	}
	summary.Title = " Summary "
	summary.Wrap = true

	// setup the keybindings of the app
	if err = g.SetKeybinding("", c.KeyCtrlN, c.ModNone, addSite); err != nil {
		log.Fatal("Failed to set keybindings")
	}
	if err = g.SetKeybinding("", c.KeyDelete, c.ModNone, deleteEntry); err != nil {
		log.Fatal("Failed to set keybindings")
	}
	if err = g.SetKeybinding(NEWS_VIEW, c.KeyCtrlB, c.ModNone, addBookmark); err != nil {
		log.Fatal("Failed to set keybindings")
	}
	if err = g.SetKeybinding("", c.KeyCtrlC, c.ModNone, quit); err != nil {
		log.Fatal("Failed to set keybindings")
	}
	if err = g.SetKeybinding("", c.KeyCtrlB, c.ModAlt, loadBookmarks); err != nil {
		log.Fatal("Failed to set keybindings")
	}
	if err = g.SetKeybinding("", c.KeyTab, c.ModNone, switchView); err != nil {
		log.Fatal("Failed to set keybindings")
	}
	if err = g.SetKeybinding("", c.KeyArrowUp, c.ModNone, listUp); err != nil {
		log.Fatal("Failed to set keybindings")
	}
	if err = g.SetKeybinding("", c.KeyArrowDown, c.ModNone, listDown); err != nil {
		log.Fatal("Failed to set keybindings")
	}
	if err = g.SetKeybinding("", c.KeyPgup, c.ModNone, listPgUp); err != nil {
		log.Fatal("Failed to set keybindings")
	}
	if err = g.SetKeybinding("", c.KeyPgdn, c.ModNone, listPgDown); err != nil {
		log.Fatal("Failed to set keybindings")
	}
	if err = g.SetKeybinding("", c.KeyEnter, c.ModNone, onEnter); err != nil {
		log.Fatal("Failed to set keybindings")
	}
	if err = g.SetKeybinding("", c.KeyCtrlQ, c.ModNone, removeTopView); err != nil {
		log.Fatal("Failed to set keybindings")
	}
	if err = g.SetKeybinding("", c.KeyCtrlF, c.ModNone, find); err != nil {
		log.Fatal("Failed to set keybindings")
	}
	if err = g.SetKeybinding(NEWS_VIEW, c.KeyCtrlO, c.ModNone, loadContent); err != nil {
		log.Fatal("Failed to set keybindings")
	}
	// run the mainloop
	if err = g.MainLoop(); err != nil && err != c.ErrQuit {
		log.Println("terminews exited unexpectedly: ", err)
	}
	log.Println("Exiting\n")

}
