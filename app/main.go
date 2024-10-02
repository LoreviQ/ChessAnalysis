package main

import (
	"github.com/LoreviQ/ChessAnalysis/app/internal/database"
	"github.com/LoreviQ/ChessAnalysis/app/internal/gui"
	"github.com/LoreviQ/ChessAnalysis/app/internal/server"
)

func startApp() error {
	// Database
	db, err := database.NewConnection(false)
	if err != nil {
		return err
	}
	defer db.Close()

	// Server
	server, _ := server.NewServer(db)
	go server.ListenAndServe()
	defer server.Close()

	// GUI
	myGUI := gui.NewGUI(2000, 1200)
	myGUI.CreateGUI()
	return nil
}

func main() {
	startApp()
}
