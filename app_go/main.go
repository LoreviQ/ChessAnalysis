package main

import "github.com/LoreviQ/ChessAnalysis/app_go/internal/database"

func main() {
	db, err := database.NewConnection(true)
	if err != nil {
		panic(err)
	}
	moves := []string{
		"1", "e4", "g6", "2", "d4", "Bg7", "3", "e5", "d6",
		"4", "exd6", "Qxd6", "5", "Nc3", "Bxd4", "6", "Be3", "Bxc3+",
		"7", "bxc3", "h5", "8", "f4", "Qxd1+", "9", "Rxd1", "Nd7",
		"10", "Bc4", "e5", "11", "fxe5", "Nxe5", "12", "Bd4", "Nxc4",
		"13", "Bxh8", "Nb6", "14", "Nf3", "Nd7", "15", "O-O", "Nh6",
		"16", "Bf6", "Nxf6", "17", "Ng5", "Nd7", "18", "Nxf7", "Nxf7",
		"19", "Rd3", "Ke7",
	}
	err = db.InsertMoves(moves, "123456")
	if err != nil {
		panic(err)
	}
	moves, err = db.GetMoves("123456")
	if err != nil {
		panic(err)
	}
	printStr := "\""
	for move := range moves {
		printStr += moves[move] + "\", \""
	}
	print(printStr[:len(printStr)-3])
	db.Close()
}
