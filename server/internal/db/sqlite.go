package db

import (
    "database/sql"
    "fmt"
    "os"
    _ "modernc.org/sqlite"
)

func sqlitePath() string {
    p := os.Getenv("SQLITE_PATH")
    if p == "" { p = "./data/picturebook.db" }
    return p
}

func Connect() (*sql.DB, error) {
    path := sqlitePath()
    dsn := fmt.Sprintf("file:%s?_pragma=foreign_keys(ON)", path)
    return sql.Open("sqlite", dsn)
}

func EnsureSchema(db *sql.DB) error {
    stmts := []string{
        `CREATE TABLE IF NOT EXISTS books (
            id TEXT PRIMARY KEY,
            title TEXT NOT NULL,
            cover_url TEXT,
            age_min INTEGER NOT NULL,
            age_max INTEGER NOT NULL,
            tags TEXT,
            popularity_score REAL NOT NULL,
            theme_keywords TEXT,
            is_editor_pick INTEGER NOT NULL DEFAULT 0,
            status TEXT NOT NULL DEFAULT 'draft'
        );`,
        `CREATE TABLE IF NOT EXISTS pages (
            id TEXT PRIMARY KEY,
            book_id TEXT NOT NULL,
            idx INTEGER NOT NULL,
            image_url TEXT,
            duration_hint INTEGER,
            FOREIGN KEY(book_id) REFERENCES books(id) ON DELETE CASCADE
        );`,
        `CREATE TABLE IF NOT EXISTS categories (
            id TEXT PRIMARY KEY,
            name TEXT NOT NULL,
            description TEXT
        );`,
    }
    for _, s := range stmts {
        if _, err := db.Exec(s); err != nil { return err }
    }
    // add column category_id if not exists
    _, _ = db.Exec(`ALTER TABLE books ADD COLUMN category_id TEXT`)
    // ensure default category exists
    _, _ = db.Exec(`INSERT OR IGNORE INTO categories(id,name,description) VALUES('default','默认分类','系统初始化默认分类')`)
    // backfill books with NULL/empty category_id
    _, _ = db.Exec(`UPDATE books SET category_id='default' WHERE category_id IS NULL OR category_id=''`)
    return nil
}