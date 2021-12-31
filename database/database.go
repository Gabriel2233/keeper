package db

import (
	"database/sql"
	"errors"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Store struct {
    db *sql.DB
}

type Folder struct {
    Id int64
    Name string
    CreatedAt time.Time
}

type Sheet struct {
    Id int64
    Name string
    Alias string
    Data string
    CreatedAt time.Time
    Folder int64
}

func NewStore(path string) (*Store, error) {
    db, err := sql.Open("sqlite3", path)    
    if err != nil {
        return nil, err
    }

    if err = db.Ping(); err != nil {
        return nil, err
    }

    createTables(db)

    store := &Store{
        db: db,
    }

    return store, nil
}

// ### BASIC OPERATIONS ON SHEETS (CREATE, FIND, REMOVE AND LIST) ###

func (s *Store) AddSheet(folder, name, alias, data string) (int64, error) {
    folderId, err := s.FindFolderIdByName(folder)
    if err != nil {
        return -1, err
    }

    res, err := execStatement(s.db, "INSERT INTO sheets(name, alias, data, folder_id) VALUES(?, ?, ?, ?)", name, alias, data, folderId)
    if err != nil {
        return -1, err
    }

    id, err := res.LastInsertId()
    if err != nil {
        return -1, err
    }

    return id, nil
}

func (s *Store) FindSheetByAlias(alias string) (Sheet, error) {
    res := s.db.QueryRow("SELECT * FROM sheets WHERE alias = ?", alias)

    var sheet Sheet
    if err := res.Scan(&sheet.Id, &sheet.Name, &sheet.Alias, &sheet.Data, &sheet.CreatedAt, &sheet.Folder); err != nil {
        if err == sql.ErrNoRows {
            return sheet, errors.New("No such sheet")
        }
        return  sheet, err
    }

    return sheet, nil
}

func (s *Store) FindSheetById(id int64) (Sheet, error) {
    var sheet Sheet

    res := s.db.QueryRow("SELECT * FROM sheets WHERE id = ?", id)
    if err := res.Scan(&sheet.Id, &sheet.Name, &sheet.Alias, &sheet.Data, &sheet.CreatedAt, &sheet.Folder); err != nil {
        if err == sql.ErrNoRows {
            return sheet, errors.New("No such sheet")
        }
        return  sheet, err
    }

    return sheet, nil
}

func (s *Store) RemoveSheetByAlias(alias string) error {
    sheet, err := s.FindSheetByAlias(alias)
    if err != nil {
        return err
    }

    _, err = execStatement(s.db, "DELETE FROM sheets WHERE id = ?", sheet.Id)
    if err != nil {
        return err
    }

    return nil
}

func (s *Store) ListSheetsInFolder(folder string) ([]Sheet, error) {
    folderId, err := s.FindFolderIdByName(folder)
    if err != nil {
        return nil, err
    }

    sheets := make([]Sheet, 0)

    res, err := s.db.Query("SELECT * FROM sheets WHERE folder_id = ?", folderId)
    if err != nil {
        return nil, err
    }

    for res.Next() {
        sheet := Sheet{}
        res.Scan(&sheet.Id, &sheet.Name, &sheet.Alias, &sheet.Data, &sheet.CreatedAt, &sheet.Folder)

        sheets = append(sheets, sheet)
    }

    return sheets, nil
}

// ### BASIC OPERATIONS ON FOLDERS (CREATE, FIND, REMOVE AND LIST) ###

func (s *Store) AddFolder(name string) (int64, error) {
    res, err := execStatement(s.db, "INSERT INTO folders(name) VALUES(?)", name)
    if err != nil {
        return -1, err
    }

    id, err := res.LastInsertId()
    if err != nil {
        return -1, err
    }

    return id, err
}

func (s *Store) FindFolderIdByName(name string) (int64, error) {
    res := s.db.QueryRow("SELECT id FROM folders WHERE name = ?;", name)

    var folderId int64
    if err := res.Scan(&folderId); err != nil {
        if err == sql.ErrNoRows {
            return -1, errors.New("No such folder")
        }
        return  -1, err
    }

    return folderId, nil
}

func (s *Store) RemoveFolderByName(name string) error {
    id, err := s.FindFolderIdByName(name)
    if err != nil {
        return err
    }

    _, err = execStatement(s.db, "DELETE FROM folders WHERE id = ?", id)
    if err != nil {
        return err
    }

    return nil
}

func (s* Store) ListFolders() ([]Folder, error) {
    folders := make([]Folder, 0)

    res, err := s.db.Query("SELECT * FROM folders")
    if err != nil {
        return nil, err
    }

    for res.Next() {
        folder := Folder{}
        res.Scan(&folder.Id, &folder.Name, &folder.CreatedAt)

        folders = append(folders, folder)
    }

    return folders, nil
}

func execStatement(db *sql.DB, query string, args ...interface{}) (sql.Result, error) {
    stmt, err := db.Prepare(query)
    if err != nil {
        return nil, err
    }

    res, err := stmt.Exec(args...)
    if err != nil {
        return nil, err
    }

    return res, nil
}

func createTables(db *sql.DB) {
    stmts := []string{
`CREATE TABLE IF NOT EXISTS folders (
id INTEGER PRIMARY KEY,name TEXT NOT NULL UNIQUE,created_at DATETIME DEFAULT CURRENT_TIMESTAMP);`,
`CREATE TABLE IF NOT EXISTS sheets (
id INTEGER PRIMARY KEY,name TEXT NOT NULL,alias TEXT NOT NULL,data TEXT NOT NULL,created_at DATETIME DEFAULT CURRENT_TIMESTAMP, 
folder_id INTEGER,FOREIGN KEY(folder_id) REFERENCES folders(id));`,
    }
        
    tx, err := db.Begin()

    for _, stmt := range stmts {
        if err != nil {
            log.Fatalf("Failed to begin transaction: %s\n", err.Error())
        }
        _, err = tx.Exec(stmt)
        if err != nil {
            tx.Rollback()
            log.Fatalf("Error on stmt: %s\n", err.Error())
        }
    }

    if err = tx.Commit(); err != nil {
        log.Fatalf("Failed to commit transaction: %s\n", err.Error())
    }
}
