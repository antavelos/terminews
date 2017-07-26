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
	"errors"
	"fmt"
	"regexp"

	"github.com/antavelos/terminews/db"
	"github.com/mmcdole/gofeed"
)

func CheckUrl(url string) (*gofeed.Feed, error) {
	fp := gofeed.NewParser()

	return fp.ParseURL(url)
}

func DownloadEvents(url string) ([]db.Event, error) {
	feed, err := CheckUrl(url)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to retrieve news from: '%v'", url))
	}

	var events []db.Event
	for _, item := range feed.Items {
		e := db.Event{}
		e.Title = item.Title
		if item.Author != nil {
			e.Author = item.Author.Name
		} else {
			e.Author = "Unknown author"
		}
		e.Url = item.Link
		e.Summary = trimHtml(item.Description)
		e.Published = item.Published

		events = append(events, e)
	}

	return events, nil
}

func trimHtml(desc string) string {
	var re = regexp.MustCompile(`(<.*?>)`)

	return re.ReplaceAllString(desc, ``)
}
