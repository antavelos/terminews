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

type Site struct {
	Id   int
	Name string
	Url  string
}

func GetSiteSql() string {
	return `
    CREATE TABLE IF NOT EXISTS site(
        Id INTEGER NOT NULL PRIMARY KEY ASC,
        Name TEXT,
        Url TEXT,
        CreatedAt DATETIME
    );`
}

func (tdb *TDB) GetSites() ([]Site, error) {
	sql_readall := `
    SELECT Id, Name, Url FROM site
    ORDER BY datetime(CreatedAt) ASC
    `

	rows, err := tdb.Query(sql_readall)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	var records []Site
	for rows.Next() {
		rr := Site{}
		if err := rows.Scan(&rr.Id, &rr.Name, &rr.Url); err != nil {
			return nil, err
		}
		records = append(records, rr)
	}
	return records, nil
}

func (tdb *TDB) GetSiteById(id int) (Site, error) {
	sql_readone := `SELECT Id, Name, Url FROM site WHERE id = ?`

	stmt, err := tdb.Prepare(sql_readone)
	defer stmt.Close()
	if err != nil {
		return Site{}, err
	}

	var rr Site
	if err = stmt.QueryRow(id).Scan(&rr.Id, &rr.Name, &rr.Url); err != nil {
		if err == sql.ErrNoRows {
			return Site{}, NotFound(fmt.Sprintf("Site not found for id: %v", id))
		}
		return Site{}, err
	}

	return rr, nil
}

func (tdb *TDB) GetSiteByUrl(url string) (Site, error) {
	sql_readone := `SELECT Id, Name, Url FROM site WHERE Url = ?`

	stmt, err := tdb.Prepare(sql_readone)
	defer stmt.Close()
	if err != nil {
		return Site{}, err
	}

	var rr Site
	if err = stmt.QueryRow(url).Scan(&rr.Id, &rr.Name, &rr.Url); err != nil {
		if err == sql.ErrNoRows {
			return Site{}, NotFound(fmt.Sprintf("Site not found for url: %v", url))
		}
		return Site{}, err
	}

	return rr, nil
}

func (tdb *TDB) AddSite(rr Site) error {
	sql_additem := `
    INSERT OR REPLACE INTO site(
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

func (tdb *TDB) DeleteSite(id int) error {
	if _, err := tdb.GetSiteById(id); err != nil {
		return NotFound(fmt.Sprintf("Site not found for id: %v", id))
	}

	sql_delete := `DELETE FROM site WHERE id = ?`

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

func (rr Site) String() string {
	return rr.Name
}
