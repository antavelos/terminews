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
	"golang.org/x/net/html"
	"net/http"
	// "strings"
)

func ParseUrl(url string, ch chan string, done chan bool) {
	resp, err := http.Get(url)
	defer func() {
		// Notify that we're done after this function
		done <- true
	}()
	if err != nil {
		return
	}
	b := resp.Body
	defer b.Close() // close Body when the function returns

	z := html.NewTokenizer(b)

	inParagraph := false
	for {
		tt := z.Next()

		switch tt {
		case html.ErrorToken:
			// End of the document, we're done
			return
		case html.StartTagToken:
			t := z.Token()

			// Check if the token is an <p> tag
			isParagraph := t.Data == "p"
			if isParagraph {
				inParagraph = true
			}
		case html.EndTagToken:
			t := z.Token()

			isParagraph := t.Data == "p"
			if isParagraph {
				inParagraph = false
			}
		case html.TextToken:
			if inParagraph {
				t := fmt.Sprint(z.Token())
				ch <- html.UnescapeString(t)
			}
		}
	}
}

func GetContent(url string) []string {
	ch := make(chan string)
	done := make(chan bool)
	go ParseUrl(url, ch, done)

	content := []string{}
	for {
		select {
		case text := <-ch:
			content = append(content, text)
		case <-done:
			return content
		}
	}
}
