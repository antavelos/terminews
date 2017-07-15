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

func GetRssReaders(db *sql.DB) ([]RssReader, error) {
	sql_readall := `
    SELECT Id, Name, Url FROM rssreader
    ORDER BY datetime(CreatedAt) ASC
    `

	rows, err := db.Query(sql_readall)
	if err != nil {
		return nil, DbError(err.Error())
	}
	defer rows.Close()

	var records []RssReader
	for rows.Next() {
		rr := RssReader{}
		if err := rows.Scan(&rr.Id, &rr.Name, &rr.Url); err != nil {
			return nil, DbError(err.Error())
		}
		records = append(records, rr)
	}
	return records, nil
}

func GetRssReaderById(db *sql.DB, id int) (RssReader, error) {
	sql_readone := `SELECT Id, Name, Url FROM rssreader WHERE id = ?`

	stmt, err := db.Prepare(sql_readone)
	if err != nil {
		return RssReader{}, DbError(err.Error())
	}
	defer stmt.Close()

	var rr RssReader
	if err = stmt.QueryRow(id).Scan(&rr.Id, &rr.Name, &rr.Url); err != nil {
		if err == sql.ErrNoRows {
			return RssReader{}, NotFound(fmt.Sprintf("RssReader not found for id: %v", id))
		}
		return RssReader{}, DbError(err.Error())
	}

	return rr, nil
}

func AddRssReader(db *sql.DB, rr RssReader) error {
	sql_additem := `
    INSERT OR REPLACE INTO rssreader(
        Name,
        Url,
        CreatedAt
    ) values(?, ?, CURRENT_TIMESTAMP)
    `

	stmt, err := db.Prepare(sql_additem)
	if err != nil {
		return DbError(err.Error())
	}
	defer stmt.Close()

	if _, err = stmt.Exec(rr.Name, rr.Url); err != nil {
		return DbError(err.Error())
	}

	return nil
}

func DeleteRssReader(db *sql.DB, id int) error {
	if _, err := GetRssReaderById(db, id); err != nil {
		return NotFound(fmt.Sprintf("RssReader not found for id: %v", id))
	}

	sql_delete := `DELETE FROM rssreader WHERE id = ?`

	stmt, err := db.Prepare(sql_delete)
	if err != nil {
		return DbError(err.Error())
	}
	defer stmt.Close()

	if _, err = stmt.Exec(id); err != nil {
		return DbError(err.Error())
	}

	return nil
}
