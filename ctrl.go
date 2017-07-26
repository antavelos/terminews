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
	"fmt"
	"strings"
	_ "time"

	"github.com/antavelos/terminews/db"
	c "github.com/jroimartin/gocui"
)

func addKeybinding(g *c.Gui, viewname string, key interface{}, mod c.Modifier, handler func(*c.Gui, *c.View) error) {
	err := g.SetKeybinding(viewname, key, mod, handler)
	if err != nil {
		fmt.Errorf("Could not set key binding: %v", err)
	}
}

func updateSummary() {
	if newsList.IsEmpty() {
		return
	}
	currItem := newsList.CurrentItem()
	event := currItem.(db.Event)

	summary.Clear()
	fmt.Fprintf(summary, "\n\n %v %v\n", Bold.Sprint("By:"), event.Author)
	fmt.Fprintf(summary, " %v %v\n", Bold.Sprint("Published on:"), event.Published)
	fmt.Fprintf(summary, " %v %v\n\n", Bold.Sprint("Site:"), event.Host())
	fmt.Fprintf(summary, " %v", event.Summary)
}

func updateNews(g *c.Gui, events []db.Event, from string) {
	if len(events) == 0 {
		newsList.SetTitle(fmt.Sprintf("No news in %v", from))
		summary.Clear()
		return
	}
	data := make([]interface{}, len(events))
	for i, e := range events {
		data[i] = e
	}

	if err := newsList.SetItems(data); err != nil {
		fmt.Errorf("Failed to update news list: %v", err)
	}
	newsList.SetTitle(fmt.Sprintf("News from: %v", from))
	newsList.Focus(g)
	newsList.ResetCursor()
	updateSummary()
}

func loadSites() {
	sites, err := tdb.GetSites()
	if err != nil {
		fmt.Errorf("Failed to load sites: %v", err)
	}
	if len(sites) == 0 {
		sitesList.SetTitle(fmt.Sprintf("No sites available"))
		sitesList.Reset()
		newsList.Reset()
		newsList.SetTitle("No news yet...")
		return
	}
	data := make([]interface{}, len(sites))
	for i, rr := range sites {
		data[i] = rr
	}

	if err := sitesList.SetItems(data); err != nil {
		fmt.Errorf("Failed to update sites list: %v", err)
	}

	sitesList.SetTitle("Sites")
}

func createPromptView(g *c.Gui, title string) error {
	tw, th := g.Size()
	v, err := g.SetView(PROMPT_VIEW, tw/6, (th/2)-1, (tw*5)/6, (th/2)+1)
	if err != nil && err != c.ErrUnknownView {
		return err
	}
	v.Editable = true
	setPromptViewTitle(g, title)
	// v.Title = title

	g.Cursor = true
	g.SetCurrentView(PROMPT_VIEW)

	return nil
}

func deletePromptView(g *c.Gui) error {
	if err := g.DeleteView(PROMPT_VIEW); err != nil {
		return err
	}
	g.Cursor = false

	return nil
}

func setPromptViewTitle(g *c.Gui, title string) {
	v, _ := g.View(PROMPT_VIEW)
	v.Title = fmt.Sprintf("%v (Ctrl-q to cancel)", title)
}

func isNewSitePrompt(v *c.View) bool {
	return strings.Contains(v.Title, "New site") || strings.Contains(v.Title, "try again")
}

func isFindPrompt(v *c.View) bool {
	return strings.Contains(v.Title, "Search ")
}

func eventSatisfiesSearch(terms []string, e db.Event) bool {
	for _, term := range terms {
		tl := strings.ToLower(term)
		title := strings.ToLower(e.Title)
		summary := strings.ToLower(e.Summary)
		if !strings.Contains(title, tl) && !strings.Contains(summary, tl) {
			return false
		}
	}
	return true
}

func findEvents(terms []string) chan db.Event {
	c := make(chan db.Event)
	sites, err := tdb.GetSites()
	if err != nil {
		close(c)
	}
	go func() {
		for _, site := range sites {
			events, err := DownloadEvents(site.Url)
			if err != nil {
				continue
			}
			for _, e := range events {
				if eventSatisfiesSearch(terms, e) {
					c <- e
				}
			}
		}
		close(c)
	}()
	return c
}

// Key binding functions

func quit(g *c.Gui, v *c.View) error {
	return c.ErrQuit
}

func switchView(g *c.Gui, v *c.View) error {
	if v == sitesList.View {
		newsList.Focus(g)
		sitesList.Unfocus()
	} else {
		sitesList.Focus(g)
		newsList.Unfocus()
	}
	return nil
}

func listUp(g *c.Gui, v *c.View) error {
	if v == sitesList.View {
		sitesList.MoveUp()
	} else {
		if !newsList.IsEmpty() {
			newsList.MoveUp()
			updateSummary()
		}
	}

	return nil
}

func listDown(g *c.Gui, v *c.View) error {
	if v == sitesList.View {
		sitesList.MoveDown()
	} else {
		if !newsList.IsEmpty() {
			newsList.MoveDown()
			updateSummary()
		}
	}

	return nil
}

func listPgDown(g *c.Gui, v *c.View) error {
	if v == sitesList.View {
		sitesList.MovePgDown()
	} else {
		if !newsList.IsEmpty() {
			newsList.MovePgDown()
			updateSummary()
		}
	}

	return nil
}

func listPgUp(g *c.Gui, v *c.View) error {
	if v == sitesList.View {
		sitesList.MovePgUp()
	} else {
		if !newsList.IsEmpty() {
			newsList.MovePgUp()
			updateSummary()
		}
	}

	return nil
}

func onEnter(g *c.Gui, v *c.View) error {
	switch v.Name() {
	case SITES_VIEW:
		currItem := sitesList.CurrentItem()
		site := currItem.(db.Site)

		summary.Clear()
		newsList.Clear()
		newsList.Focus(g)
		newsList.Title = " Downloading ... "
		g.Execute(func(g *c.Gui) error {
			events, err := DownloadEvents(site.Url)
			if err != nil {
				newsList.Title = fmt.Sprintf(" Failed to load news from: %v ", site.Name)
				newsList.Clear()
			} else {
				updateNews(g, events, site.Name)
			}
			return nil
		})
	case PROMPT_VIEW:
		if isNewSitePrompt(v) {
			url := strings.TrimSpace(v.ViewBuffer())
			if len(url) == 0 {
				return nil
			}
			g.Execute(func(g *c.Gui) error {
				feed, err := CheckUrl(url)
				if err != nil {
					setPromptViewTitle(g, "Invalid URL, try again:")
					g.SelFgColor = c.ColorRed | c.AttrBold
					return nil
				}

				_, err = tdb.GetSiteByUrl(url)
				if _, ok := err.(db.NotFound); !ok {
					setPromptViewTitle(g, "Site already exists, try again:")
					g.SelFgColor = c.ColorRed | c.AttrBold
					return nil
				}

				rr := db.Site{Name: feed.Title, Url: url}
				if err := tdb.AddSite(rr); err != nil {
					return err
				}
				deletePromptView(g)
				g.SelFgColor = c.ColorGreen | c.AttrBold
				loadSites()
				sitesList.Focus(g)

				return nil
			})
		}
		if isFindPrompt(v) {
			newsList.Reset()
			newsList.Focus(g)
			newsList.Title = " Searching ... "
			deletePromptView(g)
			terms := strings.Split(strings.TrimSpace(v.ViewBuffer()), " ")
			g.Execute(func(g *c.Gui) error {
				c := 0
				for event := range findEvents(terms) {
					newsList.AddItem(g, event)
					c++
				}
				if c == 0 {
					newsList.SetTitle("No events found")
				} else {
					newsList.SetTitle(fmt.Sprintf("%v event(s) found", c))
				}
				return nil
			})
		}
	}

	return nil
}

func addBookmark(g *c.Gui, v *c.View) error {
	switch v.Name() {
	case NEWS_VIEW:
		currItem := newsList.CurrentItem()
		event := currItem.(db.Event)
		if err := tdb.AddEvent(event); err != nil {
			return err
		}
	}
	return nil
}

func loadBookmarks(g *c.Gui, v *c.View) error {
	events, err := tdb.GetEvents()
	source := "My bookmarks"
	if err != nil {
		newsList.Title = fmt.Sprintf(" Failed to load news from: %v ", source)
		newsList.Clear()
	} else {
		updateNews(g, events, source)
	}
	return nil
}

func deleteEntry(g *c.Gui, v *c.View) error {
	if v == sitesList.View {
		currItem := sitesList.CurrentItem()
		rr := currItem.(db.Site)
		if err := tdb.DeleteSite(rr.Id); err != nil {
			return err
		}
		loadSites()
	} else {
		if strings.Contains(newsList.Title, "My bookmarks") {
			currItem := newsList.CurrentItem()
			event := currItem.(db.Event)
			if err := tdb.DeleteEvent(event.Id); err != nil {
				return err
			}
			loadBookmarks(g, v)
		}
	}
	return nil
}

func removePrompt(g *c.Gui, v *c.View) error {
	if v.Name() == PROMPT_VIEW {
		sitesList.Focus(g)
		g.SelFgColor = c.ColorGreen | c.AttrBold
		return deletePromptView(g)
	}
	return nil
}

func addSite(g *c.Gui, v *c.View) error {
	if err := createPromptView(g, "New site URL:"); err != nil {
		return err
	}

	return nil
}

func find(g *c.Gui, v *c.View) error {
	if err := createPromptView(g, "Search with multiple terms:"); err != nil {
		return err
	}

	return nil
}
