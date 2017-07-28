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
	"strings"
	"testing"
)

func TestJustifiedLines(t *testing.T) {
	text := "this is some text 1 this is some text 2 this is some text 3 this is some text"
	for _, c := range []struct {
		text string
		w    int
		want []string
	}{
		{
			text,
			10,
			[]string{
				"this is",
				"some text",
				"1 this is",
				"some text",
				"2 this is",
				"some text",
				"3 this is",
				"some text",
			},
		},
		{
			text,
			21,
			[]string{
				"this is some text 1",
				"this is some text 2",
				"this is some text 3",
				"this is some text",
			},
		},
		{
			text,
			30,
			[]string{
				"this is some text 1 this is",
				"some text 2 this is some text",
				"3 this is some text",
			},
		},
	} {
		got := JustifiedLines(c.text, c.w)
		if len(got) != len(c.want) {
			t.Errorf("for w=%d: got num of lines == %v, want %v", c.w, len(got), len(c.want))
		}
		for i := range got {
			if strings.Compare(got[i], c.want[i]) != 0 {
				t.Errorf("for w=%d: line %d got == '%v', want '%v'", c.w, i+1, got[i], c.want[i])

			}
		}
	}
}
