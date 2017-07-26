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
	"path"

	_ "github.com/mattn/go-sqlite3"
)

type TDB struct {
	*sql.DB
}

func InitDB(appDir string) (*TDB, error) {

	dbpath := path.Join(appDir, "terminews.db")
	db, err := sql.Open("sqlite3", dbpath)
	if err != nil || db == nil {
		return nil, err
	}

	tdb := &TDB{db}
	if err = tdb.CreateTables(); err != nil {
		return nil, err
	}

	return tdb, nil
}

func (tdb *TDB) CreateTables() error {
	ssql := []string{
		GetSiteSql(),
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
		"DROP TABLE site;",
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
