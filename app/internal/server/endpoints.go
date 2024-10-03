package server

import (
	"net/http"
)

// GET /readiness handler
//
// This handler is used to check if the server is ready to accept requests.
func (cfg *serverCfg) getReadiness(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

type postMovesRequest struct {
	Moves []string `json:"moves"`
}

// POST /games/{id}/moves handler
//
// This handler is used to insert moves into the database.
func (cfg *serverCfg) postMoves(w http.ResponseWriter, r *http.Request) {
	// Get game ID
	id := r.PathValue("id")
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "no ID provided")
		return
	}

	// Decode request
	var request postMovesRequest
	err := decodeRequest(w, r, &request)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "failed to decode request body")
		return
	}

	// Insert moves into database
	err = cfg.db.InsertMoves(request.Moves, id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error inserting moves into db")
		return
	}

	// Response
	w.WriteHeader(http.StatusOK)
}

type getLatestMoveResponse struct {
	Moves []string `json:"moves"`
}

// GET /games/{id}/moves/latest handler
//
// This handler is used to get the latest moves of a game from the database.
func (cfg *serverCfg) getLatestMoves(w http.ResponseWriter, r *http.Request) {
	// Get game ID
	id := r.PathValue("id")
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "no ID provided")
		return
	}

	// Get moves from database
	moves, err := cfg.db.GetMovesByChessdotcomID(id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error getting moves from db")
		return
	}

	// Response
	respondWithJSON(w, http.StatusOK, getLatestMoveResponse{Moves: moves})
}
