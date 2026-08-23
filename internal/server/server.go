// Core server: accept loop, event loop
package server

import (
	"net"
	"redis-go/internal/storage"
)

type Server struct {
	store *storage.Store
	addr  string
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", ":6379")
	if err != nil {
		return err
	}
	defer ln.Close()

	for {
		conn, _ := ln.Accept()
		go s.handle(conn)
	}
}
