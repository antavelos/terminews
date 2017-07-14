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

	result := GetRssReaders(sqldb)
	t.Log(result)
	for i, res := range result {
		if res.Name != items[i].Name {
			t.Errorf("RssReader record %v has name %v, want %v",
				res.Id, res.Name, items[i].Name)
		}
	}

	record := GetRssReaderById(sqldb, 1)
	t.Log(record)
	if record.Name != items[0].Name {
		t.Errorf("RssReader record %v has name %v, want %v",
			record.Id, record.Name, items[0].Name)
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

	result := GetBookmarks(sqldb)
	t.Log(result)
	for i, res := range result {
		if res.Title != items[i].Title {
			t.Errorf("Bookmark record %v has title %v, want %v",
				res.Id, res.Title, items[i].Title)
		}
	}

	record := GetBookmarkById(sqldb, 1)
	t.Log(record)
	if record.Title != items[0].Title {
		t.Errorf("Bookmark record %v has title %v, want %v",
			record.Id, record.Title, items[0].Title)
	}
}
