package database

import "testing"

func TestGetGames(t *testing.T) {
	// Change the working directory to the root of the project
	restore := changeDirectoryToRoot()
	defer restore()

	db, err := NewConnection(true)
	if err != nil {
		t.Error(err)
	}
	// First insert moves so games are created
	gameids := []string{"123", "456", "789"}
	for _, gameid := range gameids {
		err = db.InsertMoves([]string{"1", "e4", "e5"}, gameid)
		if err != nil {
			t.Error(err)
		}
	}
	games, err := db.GetGames()
	if err != nil {
		t.Error(err)
	}
	if len(games) != 3 {
		t.Errorf("Expected 3 games, got %d", len(games))
	}
	for i, game := range games {
		if game.ID != i+1 {
			t.Errorf("Expected game ID %d, got %d", i+1, game.ID)
		}
		if game.ChessdotcomID != gameids[i] {
			t.Errorf("Expected game chess.com ID %s, got %s", gameids[i], game.ChessdotcomID)
		}
	}

	db.Close()
}
