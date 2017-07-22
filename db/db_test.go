package db

import (
	_ "github.com/mattn/go-sqlite3"
	"os"
	"testing"
)

const dbpath = "test.db"

var tdb *TDB

func SetUp() {
	tdb, _ = InitDB(dbpath)
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

func TestRssReader(t *testing.T) {

	items := []RssReader{
		RssReader{Name: "CNN", Url: "www.cnn.com"},
		RssReader{Name: "BBC", Url: "www.bbc.com"},
	}
	for _, item := range items {
		tdb.AddRssReader(item)
	}

	result, _ := tdb.GetRssReaders()
	t.Log(result)
	if len(result) != 2 {
		t.Errorf("Found %v RssReader records %v, want %v",
			len(result), len(items))
	}
	for i, res := range result {
		if res.Name != items[i].Name {
			t.Errorf("RssReader record %v has name %v, want %v",
				res.Id, res.Name, items[i].Name)
		}
	}

	record, _ := tdb.GetRssReaderById(1)
	t.Log(record)
	if record.Name != items[0].Name {
		t.Errorf("RssReader record %v has name %v, want %v",
			record.Id, record.Name, items[0].Name)
	}

	_, err := tdb.GetRssReaderById(12345)
	if _, ok := err.(NotFound); !ok {
		t.Errorf("Expected NotFound error for id 12345")
	}

	err = tdb.DeleteRssReader(12345)
	if _, ok := err.(NotFound); !ok {
		t.Errorf("Expected NotFound error for id 12345")
	}

	tdb.DeleteRssReader(1)
	_, err = tdb.GetRssReaderById(1)
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
