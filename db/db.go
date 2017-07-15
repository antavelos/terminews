package db

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type DbError string
type NotFound string

func (e DbError) Error() string {
	return fmt.Sprintf("Generic DB error: %v", e)
}

func (e NotFound) Error() string {
	return string(e)
}

func InitDB(filepath string) *sql.DB {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		panic(err)
	}
	if db == nil {
		panic("db nil")
	}
	return db
}

func CreateTables(db *sql.DB) {
	ssql := []string{
		GetRssReaderSql(),
		GetBookmarkSql(),
	}
	for _, s := range ssql {
		_, err := db.Exec(s)
		if err != nil {
			panic(err)
		}
	}
}

func DropTables(db *sql.DB) {
	ssql := []string{
		"DROP TABLE rssreader;",
		"DROP TABLE bookmark;",
	}
	for _, s := range ssql {
		_, err := db.Exec(s)
		if err != nil {
			panic(err)
		}
	}
}
