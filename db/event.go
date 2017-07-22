package db

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type Event struct {
	Id        int
	Title     string
	Author    string
	Url       string
	Summary   string
	Published string
}

func GetEventSql() string {
	return `
    CREATE TABLE IF NOT EXISTS event(
        Id INTEGER NOT NULL PRIMARY KEY ASC,
        Title TEXT,
        Author TEXT,
        Url TEXT,
        Summary TEXT,
        Published TEXT
    );`
}

func (tdb *TDB) GetEvents() ([]Event, error) {
	sql_readall := `
    SELECT Id, Title, Author, Url, Summary, Published FROM event
    ORDER BY id ASC
    `

	rows, err := tdb.Query(sql_readall)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	var result []Event
	for rows.Next() {
		e := Event{}
		if err := rows.Scan(&e.Id, &e.Title, &e.Author, &e.Url, &e.Summary, &e.Published); err != nil {
			return nil, err
		}
		result = append(result, e)
	}
	return result, nil
}

func (tdb *TDB) GetEventById(id int) (Event, error) {
	sql_readone := `SELECT Id, Title, Author, Url, Summary, Published FROM event WHERE id = ?`

	stmt, err := tdb.Prepare(sql_readone)
	defer stmt.Close()
	if err != nil {
		return Event{}, err
	}

	var e Event
	if err = stmt.QueryRow(id).Scan(&e.Id, &e.Title, &e.Author, &e.Url, &e.Summary, &e.Published); err != nil {
		if err == sql.ErrNoRows {
			return Event{}, NotFound(fmt.Sprintf("Event not found for id: %v", id))
		}
		return Event{}, err
	}

	return e, nil
}

func (tdb *TDB) AddEvent(e Event) error {
	sql_additem := `
    INSERT OR REPLACE INTO event(
        Title,
        Author,
        Url,
        Summary,
        Published
    ) values(?, ?, ?, ?, ?)
    `

	stmt, err := tdb.Prepare(sql_additem)
	defer stmt.Close()
	if err != nil {
		return err
	}

	if _, err = stmt.Exec(e.Title, e.Author, e.Url, e.Summary, e.Published); err != nil {
		return err
	}

	return nil
}

func (tdb *TDB) DeleteEvent(id int) error {
	if _, err := tdb.GetEventById(id); err != nil {
		return NotFound(fmt.Sprintf("Event not found for id: %v", id))
	}

	sql_delete := `DELETE FROM event WHERE id = ?`

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

func (e Event) String() string {
	return string(e.Title)
}
