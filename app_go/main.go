package main

import (
	"gioui.org/app"
	"github.com/LoreviQ/ChessAnalysis/app_go/internal/database"
	"github.com/LoreviQ/ChessAnalysis/app_go/internal/gui"
	"github.com/LoreviQ/ChessAnalysis/app_go/internal/server"
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
	go gui.CreateGUI()
	app.Main()
	return nil
}

func main() {
	startApp()
}
