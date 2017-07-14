package db

import (
	"database/sql"
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

func GetRssReaders(db *sql.DB) []RssReader {
	sql_readall := `
    SELECT Id, Name, Url FROM rssreader
    ORDER BY datetime(CreatedAt) ASC
    `

	rows, err := db.Query(sql_readall)
	CheckDbError(err)
	defer rows.Close()

	var result []RssReader
	for rows.Next() {
		rr := RssReader{}
		err := rows.Scan(&rr.Id, &rr.Name, &rr.Url)
		if err != nil {
			panic(err)
		}
		result = append(result, rr)
	}
	return result
}

func GetRssReaderById(db *sql.DB, id int) RssReader {
	sql_readone := `SELECT Id, Name, Url FROM rssreader WHERE id = ?`

	stmt, err := db.Prepare(sql_readone)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	var rr RssReader
	err = stmt.QueryRow(id).Scan(&rr.Id, &rr.Name, &rr.Url)
	if err != nil {
		panic(err)
	}

	return rr
}

func AddRssReader(db *sql.DB, rr RssReader) {
	sql_additem := `
    INSERT OR REPLACE INTO rssreader(
        Name,
        Url,
        CreatedAt
    ) values(?, ?, CURRENT_TIMESTAMP)
    `

	stmt, err := db.Prepare(sql_additem)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(rr.Name, rr.Url)
	if err != nil {
		panic(err)
	}
}
