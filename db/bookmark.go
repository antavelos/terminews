package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type Bookmark struct {
	Id    int
	Title string
	Url   string
}

func GetBookmarkSql() string {
	return `
    CREATE TABLE IF NOT EXISTS bookmark(
        Id INTEGER NOT NULL PRIMARY KEY ASC,
        Title TEXT,
        Url TEXT,
        CreatedAt DATETIME
    );`
}

func GetBookmarks(db *sql.DB) []Bookmark {
	sql_readall := `
    SELECT Id, Title, Url FROM bookmark
    ORDER BY datetime(CreatedAt) ASC
    `

	rows, err := db.Query(sql_readall)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var result []Bookmark
	for rows.Next() {
		b := Bookmark{}
		err := rows.Scan(&b.Id, &b.Title, &b.Url)
		if err != nil {
			panic(err)
		}
		result = append(result, b)
	}
	return result
}

func GetBookmarkById(db *sql.DB, id int) Bookmark {
	sql_readone := `SELECT Id, Title, Url FROM bookmark WHERE id = ?`

	stmt, err := db.Prepare(sql_readone)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	var b Bookmark
	err = stmt.QueryRow(string(id)).Scan(&b.Id, &b.Title, &b.Url)
	if err != nil {
		panic(err)
	}

	return b
}

func AddBookmark(db *sql.DB, b Bookmark) {
	sql_additem := `
    INSERT OR REPLACE INTO bookmark(
        Title,
        Url,
        CreatedAt
    ) values(?, ?, CURRENT_TIMESTAMP)
    `

	stmt, err := db.Prepare(sql_additem)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(b.Title, b.Url)
	if err != nil {
		panic(err)
	}
}
