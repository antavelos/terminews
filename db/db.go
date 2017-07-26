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
