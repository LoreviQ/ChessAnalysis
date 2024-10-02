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

func TestInsertMoves(t *testing.T) {
	db, err := NewConnection(true)
	if err != nil {
		t.Error(err)
	}
	moves_to_insert := []string{
		"1", "e4", "g6", "2", "d4", "Bg7", "3", "e5", "d6",
		"4", "exd6", "Qxd6", "5", "Nc3", "Bxd4", "6", "Be3", "Bxc3+",
		"7", "bxc3", "h5", "8", "f4", "Qxd1+", "9", "Rxd1", "Nd7",
		"10", "Bc4", "e5", "11", "fxe5", "Nxe5", "12", "Bd4", "Nxc4",
		"13", "Bxh8", "Nb6", "14", "Nf3", "Nd7", "15", "O-O", "Nh6",
		"16", "Bf6", "Nxf6", "17", "Ng5", "Nd7", "18", "Nxf7", "Nxf7",
		"19", "Rd3", "Ke7",
	}
	err = db.InsertMoves(moves_to_insert, "123456")
	if err != nil {
		t.Error(err)
	}
	db.Close()
}
