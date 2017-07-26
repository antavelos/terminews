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
package db

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

const dbpath = "test.db"

var tdb *TDB

func SetUp() {
	db, err := sql.Open("sqlite3", dbpath)
	if err != nil || db == nil {
		panic("DB failed to be initialized.")
	}
	tdb = &TDB{db}
	tdb.CreateTables()
}

func TearDown() {
	tdb.DropTables()
	tdb.Close()
	os.Remove(dbpath)
}

func TestMain(m *testing.M) {
	SetUp()
	exitVal := m.Run()
	TearDown()

	os.Exit(exitVal)
}

func TestSite(t *testing.T) {

	items := []Site{
		Site{Name: "CNN", Url: "www.cnn.com"},
		Site{Name: "BBC", Url: "www.bbc.com"},
	}
	for _, item := range items {
		tdb.AddSite(item)
	}

	result, _ := tdb.GetSites()
	t.Log(result)
	if len(result) != 2 {
		t.Errorf("Found %v Site records %v, want %v",
			len(result), len(items))
	}
	for i, res := range result {
		if res.Name != items[i].Name {
			t.Errorf("Site record %v has name %v, want %v",
				res.Id, res.Name, items[i].Name)
		}
	}

	record, _ := tdb.GetSiteById(1)
	t.Log(record)
	if record.Name != items[0].Name {
		t.Errorf("Site record %v has name %v, want %v",
			record.Id, record.Name, items[0].Name)
	}

	_, err := tdb.GetSiteById(12345)
	if _, ok := err.(NotFound); !ok {
		t.Errorf("Expected NotFound error for id 12345")
	}

	err = tdb.DeleteSite(12345)
	if _, ok := err.(NotFound); !ok {
		t.Errorf("Expected NotFound error for id 12345")
	}

	tdb.DeleteSite(1)
	_, err = tdb.GetSiteById(1)
	if _, ok := err.(NotFound); !ok {
		t.Errorf("Expected NotFound error for id 12345")
	}
}

func TestEvent(t *testing.T) {

	items := []Event{
		Event{
			Title:     "some event1",
			Author:    "author1",
			Url:       "www.news.com/some-event1",
			Summary:   "summary1",
			Published: "2017"},
		Event{
			Title:     "some event2",
			Author:    "author2",
			Url:       "www.news.com/some-event2",
			Summary:   "summary2",
			Published: "2017"},
	}
	for _, item := range items {
		tdb.AddEvent(item)
	}

	result, _ := tdb.GetEvents()
	t.Log(result)
	if len(result) != 2 {
		t.Errorf("Found %v Event records %v, want %v",
			len(result), len(items))
	}
	for i, res := range result {
		if res.Title != items[i].Title {
			t.Errorf("Event record %v has title %v, want %v",
				res.Id, res.Title, items[i].Title)
		}
	}

	record, _ := tdb.GetEventById(1)
	t.Log(record)
	if record.Title != items[0].Title {
		t.Errorf("Event record %v has title %v, want %v",
			record.Id, record.Title, items[0].Title)
	}

	_, err := tdb.GetEventById(12345)
	if _, ok := err.(NotFound); !ok {
		t.Errorf("Expected NotFound error for id 12345")
	}

	err = tdb.DeleteEvent(12345)
	if _, ok := err.(NotFound); !ok {
		t.Errorf("Expected NotFound error for id 12345")
	}

	tdb.DeleteEvent(1)
	_, err = tdb.GetEventById(1)
	if _, ok := err.(NotFound); !ok {
		t.Errorf("Expected NotFound error for id 12345")
	}
}
