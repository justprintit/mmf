package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
)

type ServerConfig struct {
	Address string `yaml:",omitempty"`
	Port    uint16 `yaml:",omitempty"`
}

type Server struct {
	http.Server
	*http.ServeMux

	listener net.Listener
}

func (sc ServerConfig) NewServer() (*Server, error) {

	// Listener
	var addr string

	if len(sc.Address) > 0 {
		addr = sc.Address
	} else {
		addr = fmt.Sprintf("127.0.0.1:%v", sc.Port)
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	// Dispatcher
	mux := http.NewServeMux()

	// Server
	s := &Server{
		ServeMux: mux,
		Server: http.Server{
			Addr:    listener.Addr().String(),
			Handler: mux,
		},

		listener: listener,
	}

	log.Printf("Listening %s", s.URL().String())
	return s, nil
}

func (s *Server) URL() *url.URL {
	return &url.URL{
		Scheme: "http",
		Host:   s.Server.Addr,
		Path:   "/",
	}
}

func (s *Server) Serve() error {
	return s.Server.Serve(s.listener)
}
