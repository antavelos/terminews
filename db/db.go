package db

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
)

type TDB struct {
	*sql.DB
}

func InitDB(filepath string) (*TDB, error) {
	tdb := &TDB{}
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return nil, err
	}
	if db == nil {
		return nil, errors.New("db nil")
	}
	tdb.DB = db
	return tdb, nil
}

func (tdb *TDB) CreateTables() error {
	ssql := []string{
		GetRssReaderSql(),
		GetEventSql(),
	}
	for _, s := range ssql {
		_, err := tdb.Exec(s)
		if err != nil {
			return err
		}
	}

	return nil
}

func (tdb *TDB) DropTables() error {
	ssql := []string{
		"DROP TABLE rssreader;",
		"DROP TABLE event;",
	}
	for _, s := range ssql {
		_, err := tdb.Exec(s)
		if err != nil {
			return err
		}
	}

	return nil
}
