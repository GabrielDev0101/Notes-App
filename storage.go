package main

import (
	"database/sql"
	"time"

	_ "modernc.org/sqlite"
)

type Note struct {
	ID    int64
	Title string
	Body  string
}

type Store struct {
	conn *sql.DB
}

func (s *Store) Init() error {
	var err error

	s.conn, err = sql.Open("sqlite", "./notes.db")
	if err != nil {
		return err
	}

	createTableStmt := `
	CREATE TABLE IF NOT EXISTS notes (
		id INTEGER NOT NULL PRIMARY KEY,
		title TEXT NOT NULL,
		body TEXT NOT NULL
	);`

	_, err = s.conn.Exec(createTableStmt)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetNotes() ([]Note, error) {
	rows, err := s.conn.Query("SELECT * FROM notes")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	notes := []Note{}

	for rows.Next() {
		var note Note

		err := rows.Scan(&note.ID, &note.Title, &note.Body)
		if err != nil {
			return nil, err
		}

		notes = append(notes, note)
	}

	return notes, nil
}

func (s *Store) SaveNote(note *Note) error {
	if note.ID == 0 {
		note.ID = time.Now().UTC().UnixNano()
	}

	upsertStmt := `
	INSERT INTO notes (id, title, body)
	VALUES (?, ?, ?)
	ON CONFLICT(id) DO UPDATE SET
		title = excluded.title,
		body = excluded.body;
	`

	_, err := s.conn.Exec(
		upsertStmt,
		note.ID,
		note.Title,
		note.Body,
	)

	if err != nil {
		return err
	}

	return nil
}
