package main

import (
	"log"
	"net"
	"net/http"
	"net/url"
)

type Router interface {
	HandleFunc(string, func(http.ResponseWriter, *http.Request))
}

type Server struct {
	http.Server
	*http.ServeMux
	listener net.Listener
}

func NewServer() (*Server, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	s := &Server{
		listener: listener,
		ServeMux: mux,
		Server: http.Server{
			Addr:    listener.Addr().String(),
			Handler: mux,
		},
	}

	log.Printf("Listening %s://%s/", "http", s.Addr)
	return s, nil
}

func (s *Server) BaseURL() url.URL {
	return url.URL{
		Scheme: "http",
		Host:   s.Server.Addr,
		Path:   "/",
	}
}

func (s *Server) Serve() {
	s.Server.Serve(s.listener)
}
