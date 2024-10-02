package server

import "net/http"

func (cfg *serverCfg) getReadiness(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
