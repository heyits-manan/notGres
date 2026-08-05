package server

import (
	"context"
	"log"
	"net"
	"sync"

	"github.com/heyits-manan/notGres/internal/protocol"
)

type Config struct {
	DataDir string
	Port    string
}

type Server struct {
	cfg      Config
	ln       net.Listener
	wg       sync.WaitGroup
	mu       sync.Mutex
	shutdown bool
}

func New(cfg Config) *Server {
	return &Server{cfg: cfg}
}

func (s *Server) Serve(ln net.Listener, ready chan<- struct{}) {
	s.ln = ln
	close(ready)
	for {
		conn, err := ln.Accept()
		if err != nil {
			if s.isShuttingDown() {
				return
			}
			log.Printf("accept: %v", err)
			continue
		}
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			s.handleConn(conn)
		}()
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()
	log.Printf("connection from %s", conn.RemoteAddr())

	if err := protocol.HandleStartup(conn); err != nil {
		log.Printf("startup failed from %s: %v", conn.RemoteAddr(), err)
		return
	}

	// TODO: query loop coming up
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.mu.Lock()
	s.shutdown = true
	s.mu.Unlock()
	if s.ln != nil {
		_ = s.ln.Close()
	}
	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (s *Server) isShuttingDown() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.shutdown
}