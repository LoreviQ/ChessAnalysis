package database

type Game struct {
	ID            int
	CreatedAt     string
	ChessdotcomID string
	PlayerIsWhite bool
}

// GetGames returns all games from the database
func (d Database) GetGames() ([]Game, error) {
	var games []Game
	rows, err := d.db.Query(d.queries["GET_GAMES"])
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var game Game
		err := rows.Scan(&game.ID, &game.CreatedAt, &game.ChessdotcomID, &game.PlayerIsWhite)
		if err != nil {
			return nil, err
		}
		games = append(games, game)
	}
	return games, nil
}
