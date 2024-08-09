package main

import (
	"database/sql"
	"log"
)

var db *sql.DB

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./channels.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS channels (
			id INTEGER PRIMARY KEY UNIQUE,
			logo TEXT,
			url TEXT,
			name TEXT
		)
	`)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
}

func closeDb() {
	db.Close()
}

func getChannelById(id int) (*Channel, error) {
	var channel Channel
	err := db.QueryRow("SELECT id, logo, url, name FROM channels WHERE id = ?", id).
		Scan(&channel.ID, &channel.Logo, &channel.URL, &channel.Name)
	if err != nil {
		return nil, err
	}
	return &channel, nil
}
