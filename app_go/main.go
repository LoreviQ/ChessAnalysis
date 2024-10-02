package main

import (
	"github.com/LoreviQ/ChessAnalysis/app_go/internal/database"
)

func main() {
	d := database.NewConnection()
	d.Close()
}
