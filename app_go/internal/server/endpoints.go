package server

import (
	"net/http"
)

type postMovesRequest struct {
	GameID string   `json:"game_id"`
	Moves  []string `json:"moves"`
}

// GET /readiness handler
//
// This handler is used to check if the server is ready to accept requests.
func (cfg *serverCfg) getReadiness(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// POST /moves handler
func (cfg *serverCfg) postMoves(w http.ResponseWriter, r *http.Request) {
	var request postMovesRequest
	err := decodeRequest(w, r, &request)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "failed to decode request body")
		return
	}
	err = cfg.db.InsertMoves(request.Moves, request.GameID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error inserting moves into db")
		return
	}
	w.WriteHeader(http.StatusOK)
}
