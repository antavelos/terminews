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
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type RssReader struct {
	Id   int
	Name string
	Url  string
}

func GetRssReaderSql() string {
	return `
    CREATE TABLE IF NOT EXISTS rssreader(
        Id INTEGER NOT NULL PRIMARY KEY ASC,
        Name TEXT,
        Url TEXT,
        CreatedAt DATETIME
    );`
}

func (tdb *TDB) GetRssReaders() ([]RssReader, error) {
	sql_readall := `
    SELECT Id, Name, Url FROM rssreader
    ORDER BY datetime(CreatedAt) ASC
    `

	rows, err := tdb.Query(sql_readall)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	var records []RssReader
	for rows.Next() {
		rr := RssReader{}
		if err := rows.Scan(&rr.Id, &rr.Name, &rr.Url); err != nil {
			return nil, err
		}
		records = append(records, rr)
	}
	return records, nil
}

func (tdb *TDB) GetRssReaderById(id int) (RssReader, error) {
	sql_readone := `SELECT Id, Name, Url FROM rssreader WHERE id = ?`

	stmt, err := tdb.Prepare(sql_readone)
	defer stmt.Close()
	if err != nil {
		return RssReader{}, err
	}

	var rr RssReader
	if err = stmt.QueryRow(id).Scan(&rr.Id, &rr.Name, &rr.Url); err != nil {
		if err == sql.ErrNoRows {
			return RssReader{}, NotFound(fmt.Sprintf("RssReader not found for id: %v", id))
		}
		return RssReader{}, err
	}

	return rr, nil
}

func (tdb *TDB) GetRssReaderByUrl(url string) (RssReader, error) {
	sql_readone := `SELECT Id, Name, Url FROM rssreader WHERE Url = ?`

	stmt, err := tdb.Prepare(sql_readone)
	defer stmt.Close()
	if err != nil {
		return RssReader{}, err
	}

	var rr RssReader
	if err = stmt.QueryRow(url).Scan(&rr.Id, &rr.Name, &rr.Url); err != nil {
		if err == sql.ErrNoRows {
			return RssReader{}, NotFound(fmt.Sprintf("RssReader not found for url: %v", url))
		}
		return RssReader{}, err
	}

	return rr, nil
}

func (tdb *TDB) AddRssReader(rr RssReader) error {
	sql_additem := `
    INSERT OR REPLACE INTO rssreader(
        Name,
        Url,
        CreatedAt
    ) values(?, ?, CURRENT_TIMESTAMP)
    `

	stmt, err := tdb.Prepare(sql_additem)
	defer stmt.Close()
	if err != nil {
		return err
	}

	if _, err = stmt.Exec(rr.Name, rr.Url); err != nil {
		return err
	}

	return nil
}

func (tdb *TDB) DeleteRssReader(id int) error {
	if _, err := tdb.GetRssReaderById(id); err != nil {
		return NotFound(fmt.Sprintf("RssReader not found for id: %v", id))
	}

	sql_delete := `DELETE FROM rssreader WHERE id = ?`

	stmt, err := tdb.Prepare(sql_delete)
	defer stmt.Close()
	if err != nil {
		return err
	}

	if _, err = stmt.Exec(id); err != nil {
		return err
	}

	return nil
}

func (rr RssReader) String() string {
	return rr.Name
}
