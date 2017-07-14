package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func CheckDbError(err error) {
	if err != nil {
		panic(err)
	}
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
