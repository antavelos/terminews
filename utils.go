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
	_ "fmt"
	"strings"
)

func JustifiedLines(text string, w int) (lines []string) {
	tokens := strings.Split(text, " ")

	c := w
	temp := []string{}
	for _, t := range tokens {
		l := len(t)

		if l >= c {
			lines = append(lines, strings.Join(temp, " "))
			c = w
			temp = []string{}
		}
		temp = append(temp, t)
		// we consider also the following space
		c -= l + 1
	}
	if len(temp) > 0 {
		lines = append(lines, strings.Join(temp, " "))
	}

	return
}
