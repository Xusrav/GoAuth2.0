package app

import (
	"net/http"

	"github.com/Xusrav/GoAuth2.0/cmd/app/handlers"
	"github.com/Xusrav/GoAuth2.0/pkg/config"
)

type Server struct {
	// mux		*mux.Router
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Run(){
	h := handlers.NewHandler()
	http.HandleFunc("/", h.HandleMain)
	http.HandleFunc("/login", h.HandleGoogleLogin)
	http.HandleFunc("/redirect", h.HandleGoogleCallback)

	panic(http.ListenAndServe(config.Host + ":" + config.Port, nil))
}

