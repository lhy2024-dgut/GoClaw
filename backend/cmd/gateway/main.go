package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"my-openclaw/internal/gateway"
)

func main() {
	fmt.Println(`
╔═══════════════════════════════════════════════════════════╗
║           My-OpenClaw Gateway v1.0.0                      ║
║           (Go + React Implementation)                      ║
╚═══════════════════════════════════════════════════════════╝
`)

	port := 18789
	model := "mock-model"

	fmt.Printf("✓ Config: port=%d, model=%s\n", port, model)

	server := gateway.New(model)

	addr := fmt.Sprintf(":%d", port)
	httpServer := &http.Server{Addr: addr, Handler: server}

	go func() {
		log.Printf("Gateway starting on %s", addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	fmt.Printf(`
✓ Gateway ready!

  WebSocket:  ws://127.0.0.1:%d/ws
  HTTP:       http://127.0.0.1:%d
  Health:     http://127.0.0.1:%d/health
  Status:     http://127.0.0.1:%d/status

`, port, port, port, port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down...")
	httpServer.Close()
	log.Println("Server stopped")
}
