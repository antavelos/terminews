# terminews [![Build Status](https://travis-ci.org/antavelos/terminews.svg?branch=master)](https://travis-ci.org/antavelos/terminews)

**terminews** is a terminal application using the [GOCUI](https://github.com/jroimartin/gocui) library that allows you to manage RSS resources and display their news feed.


## Installation

### Dependencies

* [Go](https://golang.org/)
	Since proper executables are not available yet you need to have go installed in order to compile and run it.
* [Sqlite3](https://www.sqlite.org/)
	For storing RSS readers' data and bookmarking news.

### Steps

    go get github.com/antavelos/terminews
	cd $GOPATH/src/github.com/antavelos/terminews
	go build
	./terminews


## Usage

### Layout
The terminal is split in 3 different areas:
1. **RSS Readers list** which contains the list of the user's saved RSS readers.
2. **News list** which contains the news feed (list of news' titles) of the currently selected RSS reader.
3. **Summary** which contains extra information of the currently selected event.

For both lists the items are displayed paged.

### Key bindings
 Key combination | Description
---|---
<kbd>Tab</kbd>|Focuses between the RSS Readers list and the News list alternately
<kbd>Enter</kbd>|Retrieves the news feed of the currently selected RSS reader
<kbd>Ctrl</kbd><kbd>b</kbd>|Adds the currently selected event in the bookmarks list
<kbd>Ctrl</kbd><kbd>Alt</kbd><kbd>b</kbd>|Displays the bookmarked events
<kbd>Del</kbd>|Deletes the selected RSS reader of the selected bookmarked event depending on which list is currently focused
<kbd>&uarr;</kbd>|Moves to the previous list item circularly
<kbd>&darr;</kbd>|Moves to the next list item circularly
<kbd>PgUp</kbd>|Moves to the previous list page circularly
<kbd>PgDn</kbd>|Moves to the next list page circularly
<kbd>Ctrl</kbd><kbd>c</kbd>|Exits the application


## TODO
- [x] Add search functionality among the existing RSS feeds.
- [ ] Add caching
- [ ] Display full news content.
