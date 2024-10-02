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

// First tests the InsertMoves method by inserting a list of moves into the database
// and then tests the GetMoves method by retrieving the moves from the database
func TestInsertGetMoves(t *testing.T) {
	db, err := NewConnection(true)
	if err != nil {
		t.Error(err)
	}
	movesToInsert := []string{
		"1", "e4", "g6", "2", "d4", "Bg7", "3", "e5", "d6",
		"4", "exd6", "Qxd6", "5", "Nc3", "Bxd4", "6", "Be3", "Bxc3+",
		"7", "bxc3", "h5", "8", "f4", "Qxd1+", "9", "Rxd1", "Nd7",
		"10", "Bc4", "e5", "11", "fxe5", "Nxe5", "12", "Bd4", "Nxc4",
		"13", "Bxh8", "Nb6", "14", "Nf3", "Nd7", "15", "O-O", "Nh6",
		"16", "Bf6", "Nxf6", "17", "Ng5", "Nd7", "18", "Nxf7", "Nxf7",
		"19", "Rd3", "Ke7",
	}
	expectedMoves := []string{
		"e2e4", "g7g6", "d2d4", "Bf8g7", "e4e5", "d7d6", "e5xd6", "Qd8xd6",
		"Nb1c3", "Bg7xd4", "Bc1e3", "Bd4xc3+", "b2xc3", "h7h5", "f2f4",
		"Qd6xd1+", "Ra1xd1", "Nb8d7", "Bf1c4", "e7e5", "f4xe5", "Nd7xe5",
		"Be3d4", "Ne5xc4", "Bd4xh8", "Nc4b6", "Ng1f3", "Nb6d7", "O-O", "Ng8h6",
		"Bh8f6", "Nd7xf6", "Nf3g5", "Nf6d7", "Ng5xf7", "Nh6xf7", "Rd1d3", "Ke8e7",
	}
	err = db.InsertMoves(movesToInsert, "123456")
	if err != nil {
		t.Error(err)
	}
	moves, err := db.GetMoves("123456")
	if err != nil {
		t.Error(err)
	}
	if len(moves) != len(expectedMoves) {
		t.Errorf("Expected %d moves, got %d", len(expectedMoves), len(moves))
	}
	for i := range moves {
		if moves[i] != expectedMoves[i] {
			t.Errorf("Expected move %s, got %s", expectedMoves[i], moves[i])
		}
	}
	db.Close()
}
