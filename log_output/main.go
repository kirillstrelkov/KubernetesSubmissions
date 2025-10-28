package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

var randomUUID string

const iso8601Format = "2006-01-02T15:04:05.000Z"

func main() {
	randomUUID = uuid.New().String()
	fmt.Printf("Startup: Generated and stored UUID: %s\n", randomUUID)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port

	http.HandleFunc("/", statusHandler)

	fmt.Printf("HTTP server starting on port %s\n", port)
	go func() {
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	fmt.Println("Application running. Outputting the UUID every 5 seconds and serving / endpoint (Press Ctrl+C to stop)...")

	for {
		<-ticker.C

		nowUTC := time.Now().UTC()

		timestamp := nowUTC.Format(iso8601Format)

		fmt.Printf("%s: %s\n", timestamp, randomUUID)
	}
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	nowUTC := time.Now().UTC()
	timestamp := nowUTC.Format(iso8601Format)

	response := fmt.Sprintf("Timestamp: %s\nUUID: %s\n", timestamp, randomUUID)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, response)
}
