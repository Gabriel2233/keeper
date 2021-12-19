package main

import (
	"database/sql"
	"log"
	"time"

    _ "github.com/mattn/go-sqlite3"
)

type Store struct {
    db *sql.DB
}

type Folder struct {
    id int64
    name string
    alias string
    createdAt time.Time
}

type Content struct {
    id int64
    title string
    data string
    createdAt time.Time
    folder int64
}

func NewStore(path string) (*Store, error) {
    var store *Store
    db, err := sql.Open("sqlite3", path)    
    if err != nil {
        return store, err
    }

    if err = db.Ping(); err != nil {
        return store, err
    }

    stmts := []string{
`CREATE TABLE IF NOT EXISTS folders (
id INTEGER PRIMARY KEY AUTOINCREMENT,name TEXT NOT NULL UNIQUE,alias TEXT,created_at DATETIME NOT_NULL,contents BLOB);`,
`CREATE TABLE IF NOT EXISTS contents (
id INTEGER PRIMARY KEY AUTOINCREMENT,title TEXT NOT NULL,alias TEXT UNIQUE,created_at DATETIME NOT_NULL, folder_id INTEGER,
FOREIGN KEY(folder_id) REFERENCES folders(id));`,
    }

    for _, stmt := range stmts {
        _, err := db.Exec(stmt)
        if err != nil {
            log.Fatalf("Error on stmt: %s\n", err.Error())
        }
    }
    
    store.db = db

    return store, nil
}

func (s *Store) DoSmt() {
}
