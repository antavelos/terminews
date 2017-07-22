# terminews [![Build Status](https://travis-ci.org/antavelos/terminews.svg?branch=master)](https://travis-ci.org/antavelos/terminews)

**terminews** is a terminal application using the [GOCUI](https://github.com/jroimartin/gocui) library that allows you to manage RSS resources and display their news feed

## Installation
    go get github.com/antavelos/terminews

## Usage

### Layout
The terminal is split in 3 different areas: 
1. **RSS Readers list** which contains the list of the user's saved RSS readers.
2. **News list** which contains the news feed (list of news' titles) of the currently selected RSS reader.
3. **Summary** which contains extra information of the currently selected event.

For both lists the items are displayed paged. 

### Key bindings
| Key | Description |
|-----|-------------|
|<kbd>Tab</kbd>|focuses between the RSS Readers list and the News list alternately
|<kbd>Enter</kbd>|retrieves the news feed of the currently selected RSS reader
|<kbd>b</kbd>|adds the currently selected event in the bookmarks list
|<kbd>Ctrl</kbd><kbd>b</kbd>|displayes the bookmarked events
|<kbd>d</kbd>|deletes the selected RSS reader of the selected bookmarked event depending on which list is currently focused.
|<kbd>&uarr;</kbd>|moves to the previous list item circularly.
|<kbd>&darr;</kbd>|moves to the next list item circularly.
|<kbd>PgUp</kbd>|moves to the previous list page circularly.
|<kbd>PgDn</kbd>|moves to the next list page circularly.
|<kbd>q</kbd>|exit the application


## TODO
- [ ] Add search functionality among the existing RSS readers
- [ ] Display full news content
