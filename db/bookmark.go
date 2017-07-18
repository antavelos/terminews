package db

import (
	"database/sql"
	"fmt"
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

func (tdb *TDB) GetBookmarks() ([]Bookmark, error) {
	sql_readall := `
    SELECT Id, Title, Url FROM bookmark
    ORDER BY datetime(CreatedAt) ASC
    `

	rows, err := tdb.Query(sql_readall)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	var result []Bookmark
	for rows.Next() {
		b := Bookmark{}
		if err := rows.Scan(&b.Id, &b.Title, &b.Url); err != nil {
			return nil, err
		}
		result = append(result, b)
	}
	return result, nil
}

func (tdb *TDB) GetBookmarkById(id int) (Bookmark, error) {
	sql_readone := `SELECT Id, Title, Url FROM bookmark WHERE id = ?`

	stmt, err := tdb.Prepare(sql_readone)
	defer stmt.Close()
	if err != nil {
		return Bookmark{}, err
	}

	var b Bookmark
	if err = stmt.QueryRow(id).Scan(&b.Id, &b.Title, &b.Url); err != nil {
		if err == sql.ErrNoRows {
			return Bookmark{}, NotFound(fmt.Sprintf("Bookmark not found for id: %v", id))
		}
		return Bookmark{}, err
	}

	return b, nil
}

func (tdb *TDB) AddBookmark(b Bookmark) error {
	sql_additem := `
    INSERT OR REPLACE INTO bookmark(
        Title,
        Url,
        CreatedAt
    ) values(?, ?, CURRENT_TIMESTAMP)
    `

	stmt, err := tdb.Prepare(sql_additem)
	defer stmt.Close()
	if err != nil {
		return err
	}

	if _, err = stmt.Exec(b.Title, b.Url); err != nil {
		return err
	}

	return nil
}

func (tdb *TDB) DeleteBookmark(id int) error {
	if _, err := tdb.GetBookmarkById(id); err != nil {
		return NotFound(fmt.Sprintf("Bookmark not found for id: %v", id))
	}

	sql_delete := `DELETE FROM bookmark WHERE id = ?`

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
