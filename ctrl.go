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
	"log"
	"os/exec"
	"strings"

	"github.com/antavelos/terminews/db"
	c "github.com/jroimartin/gocui"
)

// UpdateSummary updates the summary View based on the currently selected
// news item
func UpdateSummary() error {
	Summary.Clear()

	currItem := NewsList.CurrentItem()
	if currItem == nil {
		return nil
	}
	event := currItem.(db.Event)

	authorLine := fmt.Sprintf("%v %v", Bold.Sprint("By:"), event.Author)
	publishedLine := fmt.Sprintf("%v %v", Bold.Sprint("Published on:"), event.Published)
	urlLine := fmt.Sprintf("%v %v", Bold.Sprint("URL:"), event.Url)

	w, _ := Summary.Size()
	summaryLine := strings.Join(JustifiedLines(event.Summary, w-2), "\n ")

	_, err := fmt.Fprintf(Summary, "\n\n %v\n %v\n %v\n\n\n %v",
		authorLine, publishedLine, urlLine, Bold.Sprint(summaryLine))

	return err
}

func eventInBookmarks(event db.Event) (db.Event, bool) {
	for _, b := range CurrentBookmarks {
		if event.Url == b.Url {
			return b, true
		}
	}
	return db.Event{}, false
}

// UpdateNews updates the news list according to the given events
func UpdateNews(events []db.Event, from string) error {
	NewsList.Reset()
	Summary.Clear()

	if len(events) == 0 {
		NewsList.SetTitle(fmt.Sprintf("No news in %v", from))
		return nil
	}
	NewsList.SetTitle(fmt.Sprintf("News from: %v", from))

	data := make([]interface{}, len(events))
	for i, e := range events {
		if _, ok := eventInBookmarks(e); ok {
			e.Title = fmt.Sprintf("  %v", e.Title)
		}
		data[i] = e
	}

	return NewsList.SetItems(data)
}

// LoadSites loads the sites from DB and displays them in the list
func LoadSites() error {
	SitesList.SetTitle("Sites")

	sites, err := tdb.GetSites()
	if err != nil {
		fmt.Errorf("Failed to load sites: %v", err)
	}
	if len(sites) == 0 {
		SitesList.SetTitle("No sites yet... (Ctrl-n to add)")
		SitesList.Reset()
		NewsList.Reset()
		NewsList.SetTitle("No news yet...")
		return nil
	}
	data := make([]interface{}, len(sites))
	for i, rr := range sites {
		data[i] = rr
	}

	return SitesList.SetItems(data)
}

// createContentView creates a view where the contents of thecurrently selected
// event will be displayed
func createContentView(g *c.Gui) error {
	tw, th := g.Size()
	v, err := g.SetView(CONTENT_VIEW, tw/8, th/8, (tw*7)/8, (th*7)/8)
	if err != nil && err != c.ErrUnknownView {
		return err
	}
	ContentList = CreateList(v, false)
	setTopWindowTitle(g, CONTENT_VIEW, "")
	_, err = g.SetCurrentView(CONTENT_VIEW)

	return err
}

// createPromptView creates a general purpose view to be used as input source
// from the user
func createPromptView(g *c.Gui, title string) error {
	tw, th := g.Size()
	v, err := g.SetView(PROMPT_VIEW, tw/6, (th/2)-1, (tw*5)/6, (th/2)+1)
	if err != nil && err != c.ErrUnknownView {
		return err
	}
	v.Editable = true
	setTopWindowTitle(g, PROMPT_VIEW, title)

	g.Cursor = true
	_, err = g.SetCurrentView(PROMPT_VIEW)

	return err
}

// deleteContentView deletes the current prompt view
func deleteContentView(g *c.Gui) error {
	g.Cursor = false
	return g.DeleteView(CONTENT_VIEW)
}

// deletePromptView deletes the current prompt view
func deletePromptView(g *c.Gui) error {
	g.Cursor = false
	return g.DeleteView(PROMPT_VIEW)
}

func setTopWindowTitle(g *c.Gui, view_name, title string) {
	v, err := g.View(view_name)
	if err != nil {
		log.Println("Error on setTopWindowTitle", err)
		return
	}
	v.Title = fmt.Sprintf("%v (Ctrl-q to close)", title)
}

func isNewSitePrompt(v *c.View) bool {
	return strings.Contains(v.Title, "New site") || strings.Contains(v.Title, "try again")
}

func isFindPrompt(v *c.View) bool {
	return strings.Contains(v.Title, "Search ")
}

func isBookmarksNews() bool {
	return strings.Contains(NewsList.Title, "My bookmarks")
}

// eventSatisfiesSearch searches within thr title and the summary of an event
// and if a list of terms exists conjuctively and case insensitively
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

// findEvents downloads the events of every available site and returns those
// which match the given terms
func findEvents(terms []string, c chan db.Event, done chan bool) {
	defer func() {
		done <- true
	}()

	sites, err := tdb.GetSites()
	if err != nil {
		return
	}

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
}

// Key binding functions

func Quit(g *c.Gui, v *c.View) error {
	return c.ErrQuit
}

func SwitchView(g *c.Gui, v *c.View) error {
	switch v.Name() {
	case SITES_VIEW:
		g.SelFgColor = c.ColorGreen | c.AttrBold
		if v == SitesList.View {
			NewsList.Focus(g)
			SitesList.Unfocus()
			if strings.Contains(NewsList.Title, "bookmarks") {
				g.SelFgColor = c.ColorMagenta | c.AttrBold
			}
		}
	case NEWS_VIEW:
		SitesList.Focus(g)
		NewsList.Unfocus()
	}

	return nil
}

func ListUp(g *c.Gui, v *c.View) error {
	switch v.Name() {

	case SITES_VIEW:
		if err := SitesList.MoveUp(); err != nil {
			log.Println("Error on SitesList.MoveUp()", err)
			return err
		}
	case NEWS_VIEW:
		if err := NewsList.MoveUp(); err != nil {
			log.Println("Error on NewsList.MoveUp()", err)
			return err
		}
		if err := UpdateSummary(); err != nil {
			log.Println("Error on UpdateSummary()", err)
			return err
		}
	case CONTENT_VIEW:
		if err := ContentList.MoveUp(); err != nil {
			log.Println("Error on ContentList.MoveUp()", err)
			return err
		}
	}
	return nil
}

func ListDown(g *c.Gui, v *c.View) error {
	switch v.Name() {

	case SITES_VIEW:
		if err := SitesList.MoveDown(); err != nil {
			log.Println("Error on SitesList.MoveDown()", err)
			return err
		}
	case NEWS_VIEW:
		if err := NewsList.MoveDown(); err != nil {
			log.Println("Error on NewsList.MoveDown()", err)
			return err
		}
		if err := UpdateSummary(); err != nil {
			log.Println("Error on UpdateSummary()", err)
			return err
		}
	case CONTENT_VIEW:
		if err := ContentList.MoveDown(); err != nil {
			log.Println("Error on ContentList.MoveDown()", err)
			return err
		}
	}
	return nil
}

func ListPgDown(g *c.Gui, v *c.View) error {
	switch v.Name() {

	case SITES_VIEW:
		if err := SitesList.MovePgDown(); err != nil {
			log.Println("Error on SitesList.MovePgDown()", err)
			return err
		}
	case NEWS_VIEW:
		if err := NewsList.MovePgDown(); err != nil {
			log.Println("Error on NewsList.MovePgDown()", err)
			return err
		}
		if err := UpdateSummary(); err != nil {
			log.Println("Error on UpdateSummary()", err)
			return err
		}
	case CONTENT_VIEW:
		if err := ContentList.MovePgDown(); err != nil {
			log.Println("Error on ContentList.MovePgDown()", err)
			return err
		}
	}
	return nil
}

func ListPgUp(g *c.Gui, v *c.View) error {
	switch v.Name() {

	case SITES_VIEW:
		if err := SitesList.MovePgUp(); err != nil {
			log.Println("Error on SitesList.MovePgUp()", err)
			return err
		}
	case NEWS_VIEW:
		if err := NewsList.MovePgUp(); err != nil {
			log.Println("Error on NewsList.MovePgUp()", err)
			return err
		}
		if err := UpdateSummary(); err != nil {
			log.Println("Error on UpdateSummary()", err)
			return err
		}
	case CONTENT_VIEW:
		if err := ContentList.MovePgUp(); err != nil {
			log.Println("Error on ContentList.MovePgUp()", err)
			return err
		}
	}
	return nil
}

func OnEnter(g *c.Gui, v *c.View) error {
	switch v.Name() {
	case SITES_VIEW:
		currItem := SitesList.CurrentItem()
		if currItem == nil {
			return nil
		}
		site := currItem.(db.Site)

		Summary.Clear()
		NewsList.Clear()
		NewsList.Focus(g)
		g.SelFgColor = c.ColorGreen | c.AttrBold
		NewsList.Title = " Fetching ... "
		g.Update(func(g *c.Gui) error {
			events, err := DownloadEvents(site.Url)
			if err != nil {
				NewsList.Title = fmt.Sprintf(" Failed to load news from: %v ", site.Name)
				NewsList.Clear()
			} else {
				NewsList.Focus(g)
				if err := UpdateNews(events, site.Name); err != nil {
					log.Println("Error on UpdateNews", err)
					return err
				}
				if err := UpdateSummary(); err != nil {
					log.Println("Error on UpdateSummary", err)
					return err
				}
			}
			return nil
		})
	case PROMPT_VIEW:
		if isNewSitePrompt(v) {
			url := strings.TrimSpace(v.ViewBuffer())
			if len(url) == 0 {
				return nil
			}
			g.Update(func(g *c.Gui) error {
				feed, err := CheckUrl(url)
				if err != nil {
					setTopWindowTitle(g, PROMPT_VIEW, "Invalid URL, try again:")
					g.SelFgColor = c.ColorRed | c.AttrBold
					return nil
				}

				_, err = tdb.GetSiteByUrl(url)
				if err != nil {
					if _, ok := err.(db.NotFound); !ok {
						setTopWindowTitle(g, PROMPT_VIEW, "Site already exists, try again:")
						g.SelFgColor = c.ColorRed | c.AttrBold
						return nil
					}
				} else {
					log.Println("Error o GetSiteByUrl", err)
					return err
				}

				rr := db.Site{Name: feed.Title, Url: url}
				if err := tdb.AddSite(rr); err != nil {
					log.Println("Error on AddSite", err)
					return err
				}
				deletePromptView(g)
				g.SelFgColor = c.ColorGreen | c.AttrBold
				SitesList.Focus(g)

				if err = LoadSites(); err != nil {
					log.Println("Error on LoadSites", err)
					return err
				}

				return nil
			})
		}
		if isFindPrompt(v) {
			NewsList.Reset()
			NewsList.Focus(g)
			SitesList.Unfocus()
			NewsList.Title = " Searching ... "
			deletePromptView(g)
			terms := strings.Split(strings.TrimSpace(v.ViewBuffer()), " ")
			done := make(chan bool)
			cevent := make(chan db.Event)
			go findEvents(terms, cevent, done)
			go func() {
				ct := 0
				for {
					select {
					case <-done:
						g.Update(func(g *c.Gui) error {
							NewsList.SetTitle(fmt.Sprintf("%v event(s) found", ct))
							return nil
						})
						return
					case event := <-cevent:
						g.Update(func(g *c.Gui) error {
							NewsList.AddItem(g, event)
							NewsList.SetTitle(fmt.Sprintf("%v event(s) found so far...", ct))
							return nil
						})
						ct++
					}
				}
			}()
		}
	}

	return nil
}

func AddBookmark(g *c.Gui, v *c.View) error {
	var err error
	if v.Name() == NEWS_VIEW {
		g.Update(func(g *c.Gui) error {
			currItem := NewsList.CurrentItem()
			if currItem == nil {
				return nil
			}
			event := currItem.(db.Event)

			if bookmark, ok := eventInBookmarks(event); ok {
				if err := tdb.DeleteEvent(bookmark.Id); err != nil {
					log.Println("Error on DeleteEvent", err)
					return err
				}
				event.Title = event.Title[5:]
			} else {
				if err := tdb.AddEvent(event); err != nil {
					log.Println("Error on AddEvent", err)
					return err
				}
				event.Title = fmt.Sprintf("  %v", event.Title)
			}
			if CurrentBookmarks, err = tdb.GetEvents(); err != nil {
				log.Println("Error on GetEvents", err)
				return err
			}
			NewsList.UpdateCurrentItem(event)
			if err := NewsList.DrawCurrentPage(); err != nil {
				log.Println("Error while updating event on bookmark", err)
				return err
			}
			return nil
		})
	}
	return nil
}

func LoadBookmarks(g *c.Gui, v *c.View) error {
	var err error
	name := v.Name()
	if name == PROMPT_VIEW || name == CONTENT_VIEW {
		return nil
	}

	CurrentBookmarks, err = tdb.GetEvents()
	if err != nil {
		log.Println("Error on AddEvent", err)
		return err
	}
	source := "My bookmarks"
	if err != nil {
		NewsList.Title = fmt.Sprintf(" Failed to load news from: %v ", source)
		NewsList.Clear()
	} else {
		NewsList.Focus(g)
		if err := UpdateNews(CurrentBookmarks, source); err != nil {
			log.Println("Error on UpdateNews", err)
			return err
		}
		if err := UpdateSummary(); err != nil {
			log.Println("Error on UpdateSummary", err)
			return err
		}
	}
	g.SelFgColor = c.ColorMagenta | c.AttrBold
	return nil
}

func DeleteEntry(g *c.Gui, v *c.View) error {
	switch v.Name() {

	case SITES_VIEW:
		currItem := SitesList.CurrentItem()
		if currItem == nil {
			return nil
		}
		rr := currItem.(db.Site)
		if err := tdb.DeleteSite(rr.Id); err != nil {
			log.Println("Error on DeleteSite", err)
			return err
		}
		if err := LoadSites(); err != nil {
			log.Println("Error on LoadSites", err)
			return err
		}
	case NEWS_VIEW:
		if strings.Contains(NewsList.Title, "My bookmarks") {
			currItem := NewsList.CurrentItem()
			if currItem == nil {
				return nil
			}
			event := currItem.(db.Event)
			if err := tdb.DeleteEvent(event.Id); err != nil {
				log.Println("Error on DeleteEvent", err)
				return err
			}
			if err := LoadBookmarks(g, v); err != nil {
				log.Println("Error on LoadBookmarks", err)
				return err
			}
		}
	}
	return nil
}

func RemoveTopView(g *c.Gui, v *c.View) error {
	switch v.Name() {

	case PROMPT_VIEW:
		SitesList.Focus(g)
		if isBookmarksNews() {
			g.SelFgColor = c.ColorMagenta | c.AttrBold
		} else {
			g.SelFgColor = c.ColorGreen | c.AttrBold
		}
		if err := deletePromptView(g); err != nil {
			log.Println("Error on deletePromptView", err)
			return err
		}
	case CONTENT_VIEW:
		NewsList.Focus(g)
		if isBookmarksNews() {
			g.SelFgColor = c.ColorMagenta | c.AttrBold
		} else {
			g.SelFgColor = c.ColorGreen | c.AttrBold
		}
		if err := deleteContentView(g); err != nil {
			log.Println("Error on deleteContentView", err)
			return err
		}
	}
	return nil
}

func AddSite(g *c.Gui, v *c.View) error {
	if err := createPromptView(g, "New site URL:"); err != nil {
		log.Println("Error on createPromptView", err)
		return err
	}

	return nil
}

func Find(g *c.Gui, v *c.View) error {
	if err := createPromptView(g, "Search with multiple terms:"); err != nil {
		log.Println("Error on createPromptView", err)
		return err
	}

	return nil
}

func LoadContent(g *c.Gui, v *c.View) error {
	if v.Name() == NEWS_VIEW {
		if err := createContentView(g); err != nil {
			log.Println("Error on createContentView", err)
			return err
		}
		g.SelFgColor = c.ColorGreen | c.AttrBold
		cv, _ := g.View(CONTENT_VIEW)
		cv.Title = "Fetching..."
		g.Update(func(g *c.Gui) error {
			ContentList.Focus(g)
			currItem := NewsList.CurrentItem()
			if currItem == nil {
				return nil
			}
			event := currItem.(db.Event)
			CurrentContent, _ = GetContent(event.Url)
			if err := UpdateContent(g, CurrentContent); err != nil {
				log.Println("Error on UpdateContent", err)
				return err
			}
			ContentList.SetTitle(fmt.Sprintf("%v (Ctrl-q to close)", event.Title))

			return nil
		})

	}
	return nil
}

func UpdateContent(g *c.Gui, content []string) error {
	w, _ := ContentList.Size()
	ContentList.AddItem(g, "")
	for _, text := range content {
		lines := JustifiedLines(text, w-2)
		for _, l := range lines {
			err := ContentList.AddItem(g, l)
			if err != nil {
				log.Println("Error on ContentList.AddItem", err)
				return err
			}
		}
		ContentList.AddItem(g, "")
	}
	return nil
}

func OpenBrowser(g *c.Gui, v *c.View) error {
	currItem := NewsList.CurrentItem()
	if currItem == nil {
		return nil
	}
	event := currItem.(db.Event)
	if v.Name() == NEWS_VIEW {
		cmd := exec.Command("xdg-open", event.Url)

		if err := cmd.Run(); err != nil {
			log.Println("Error on opening browser", err)
			return err
		}
	}
	return nil
}
