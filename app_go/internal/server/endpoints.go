package server

import "net/http"

// GET /readiness handler
//
// This handler is used to check if the server is ready to accept requests.
func (cfg *serverCfg) getReadiness(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
