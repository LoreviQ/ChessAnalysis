package database

import "testing"

func TestNewConnection(t *testing.T) {
	db, err := NewConnection(true)
	if err != nil {
		t.Error(err)
	}
	if db.db == nil {
		t.Error("Database connection is nil")
	}
	db.Close()
}
