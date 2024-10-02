package server

import "net/http"

type serverCfg struct {
	port string
}

func newServer() *http.Server {
	return &http.Server{}
}
