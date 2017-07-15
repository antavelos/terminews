package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"testing"
)

const dbpath = "test.db"

var sqldb *sql.DB

func SetUp() {
	sqldb = InitDB(dbpath)
	CreateTables(sqldb)
}

func TearDown() {
	DropTables(sqldb)
	sqldb.Close()
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
		AddRssReader(sqldb, item)
	}

	result, _ := GetRssReaders(sqldb)
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

	record, _ := GetRssReaderById(sqldb, 1)
	t.Log(record)
	if record.Name != items[0].Name {
		t.Errorf("RssReader record %v has name %v, want %v",
			record.Id, record.Name, items[0].Name)
	}

	_, err := GetRssReaderById(sqldb, 12345)
	if _, ok := err.(NotFound); !ok {
		t.Errorf("Expected NotFound error for id 12345")
	}

	err = DeleteRssReader(sqldb, 12345)
	if _, ok := err.(NotFound); !ok {
		t.Errorf("Expected NotFound error for id 12345")
	}

	DeleteRssReader(sqldb, 1)
	_, err = GetRssReaderById(sqldb, 1)
	if _, ok := err.(NotFound); !ok {
		t.Errorf("Expected NotFound error for id 12345")
	}
}

func TestBookmark(t *testing.T) {

	items := []Bookmark{
		Bookmark{Title: "some event1", Url: "www.news.com/some-event1"},
		Bookmark{Title: "some event2", Url: "www.news.com/some-event2"},
	}
	for _, item := range items {
		AddBookmark(sqldb, item)
	}

	result, _ := GetBookmarks(sqldb)
	t.Log(result)
	if len(result) != 2 {
		t.Errorf("Found %v Bookmark records %v, want %v",
			len(result), len(items))
	}
	for i, res := range result {
		if res.Title != items[i].Title {
			t.Errorf("Bookmark record %v has title %v, want %v",
				res.Id, res.Title, items[i].Title)
		}
	}

	record, _ := GetBookmarkById(sqldb, 1)
	t.Log(record)
	if record.Title != items[0].Title {
		t.Errorf("Bookmark record %v has title %v, want %v",
			record.Id, record.Title, items[0].Title)
	}

	_, err := GetBookmarkById(sqldb, 12345)
	if _, ok := err.(NotFound); !ok {
		t.Errorf("Expected NotFound error for id 12345")
	}

	err = DeleteBookmark(sqldb, 12345)
	if _, ok := err.(NotFound); !ok {
		t.Errorf("Expected NotFound error for id 12345")
	}

	DeleteBookmark(sqldb, 1)
	_, err = GetBookmarkById(sqldb, 1)
	if _, ok := err.(NotFound); !ok {
		t.Errorf("Expected NotFound error for id 12345")
	}
}
