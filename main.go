package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/heyits-manan/notGres/internal/server"
)

func main(){
	var cfg server.Config
	flag.StringVar(&cfg.DataDir, "data-dir", "./data", "directory for database files")
	flag.StringVar(&cfg.Port, "port", "5432", "port to listen on")
	flag.Parse()

	if err := os.MkdirAll(cfg.DataDir, 0o755); err != nil {
		log.Fatalf("create data dir: %v", err)
	}

	addr := ":" + cfg.Port
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("listen on %s: %v", addr, err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	srv := server.New(cfg)
	ready := make(chan struct{})
	go srv.Serve(ln, ready)
	<-ready
	log.Printf("notgres listening on %s", ln.Addr())

	<-ctx.Done()
	log.Println("shutting down")

	sctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(sctx); err != nil {
		log.Fatalf("shutdown: %v", err)
	}
}