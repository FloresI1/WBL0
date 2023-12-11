package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	_ "github.com/lib/pq"
)

var (
	db         *sql.DB
	cache      sync.Map
	cacheMutex sync.Mutex
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go consumeDeliveries(ctx)
	http.HandleFunc("/delivery", getDeliveryAndPaymentByOrderUIDHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		renderDeliveryAndPaymentPage(w, nil)
	})

	http.HandleFunc("/debug/cache", debugCacheHandler)

	log.Println("Server is running on :8080")
	server := &http.Server{Addr: ":8080"}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	<-signalCh

	log.Println("Shutting down HTTP server...")
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}
}
