package main

import (
	"github.com/LoreviQ/ChessAnalysis/app/internal/database"
	"github.com/LoreviQ/ChessAnalysis/app/internal/gui"
	"github.com/LoreviQ/ChessAnalysis/app/internal/server"
)

func startApp() error {
	// Database
	db, err := database.NewConnection(0)
	if err != nil {
		return err
	}
	defer db.Close()

	// Server
	server, _ := server.NewServer(db)
	go server.ListenAndServe()
	defer server.Close()

	// GUI
	myGUI := gui.NewGUI(2400, 1200, db)
	myGUI.CreateGUI()
	return nil
}

func main() {
	startApp()
}
