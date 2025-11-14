package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
)

var randomUUID string

const iso8601Format = "2006-01-02T15:04:05.000Z"

func getCounter() string {
	svc := os.Getenv("PING_PONG_SERVICE")
	if svc == "" {
		svc = "localhost:8080"
	}
	url := fmt.Sprintf("http://%s/pings", svc)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Sprintf("Error making request: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Sprintf("Received non-OK status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintf("Error reading response body: %v", err)
	}

	counter, err := strconv.Atoi(string(body))

	if err != nil {
		return fmt.Sprintf("Error reading response body: %v", err)
	}

	return fmt.Sprintf("%d", counter)
}

func printConfigValues() {
	msg := os.Getenv("MESSAGE")

	fs, err := os.Open("/tmp/information.txt")
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer fs.Close()

	content, err := io.ReadAll(fs)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	fmt.Printf("file content: %s", content)
	fmt.Printf("env variable: MESSAGE=%s\n", msg)
}

func main() {
	printConfigValues()

	randomUUID = uuid.New().String()
	fmt.Printf("Startup: Generated and stored UUID: %s\n", randomUUID)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port

	http.HandleFunc("/", statusHandler)

	fmt.Printf("HTTP server starting on port %s\n", port)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("HTTP server failed: %v", err)
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

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	nowUTC := time.Now().UTC()
	timestamp := nowUTC.Format(iso8601Format)

	line := fmt.Sprintf("%s: %s\n", timestamp, randomUUID)
	counter := getCounter()
	// print
	fmt.Print(line)
	fmt.Printf("Ping / Pongs: %s\n", counter)

	// html out
	fmt.Fprintln(w, line)
	fmt.Fprintf(w, "Ping / Pongs: %s", counter)
}
