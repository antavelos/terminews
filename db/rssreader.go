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
