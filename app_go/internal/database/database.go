package database

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"runtime"

	_ "github.com/mattn/go-sqlite3"
)

func NewConnection() *sql.DB {
	// Get file paths
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("Unable to get the current file path")
	}
	dir := filepath.Dir(filename)
	databasePath := filepath.Join(dir, "database.db")
	schemaPath := filepath.Join(dir, "sql", "schema.sql")

	// Open a connection to the SQLite3 database
	db, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		panic(err)
	}

	// Read the schema.sql file
	schema, err := os.ReadFile(schemaPath)
	if err != nil {
		log.Fatal(err)
	}

	// Execute the SQL commands in the schema.sql file
	_, err = db.Exec(string(schema))
	if err != nil {
		log.Fatal(err)
	}
	return db
}
