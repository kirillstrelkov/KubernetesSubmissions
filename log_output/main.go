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

func getCounter() (string, error) {
	svc := os.Getenv("PING_PONG_SERVICE")
	if svc == "" {
		svc = "localhost:8080"
	}
	url := fmt.Sprintf("http://%s/pings", svc)

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("error making request: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-OK status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	counter, err := strconv.Atoi(string(body))

	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	return fmt.Sprintf("%d", counter), nil
}
func printConfigValues(w io.Writer) {
	fs, err := os.Open("/tmp/information.txt")
	if err != nil {
		fmt.Fprintf(w, "Error opening file: %v\n", err)
		return
	}
	defer fs.Close()

	content, err := io.ReadAll(fs)
	if err != nil {
		fmt.Fprintf(w, "Error reading file: %v\n", err)
		return
	}

	msg := os.Getenv("MESSAGE")
	fmt.Fprintf(w, "file content: %s\n", content)
	fmt.Fprintf(w, "env variable: MESSAGE=%s\n", msg)
}

func main() {
	printConfigValues(os.Stdout)

	randomUUID = uuid.New().String()
	fmt.Printf("Startup: Generated and stored UUID: %s\n", randomUUID)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port

	http.HandleFunc("/", statusHandler)
	http.HandleFunc("/healthz", statusAlive)

	fmt.Printf("HTTP server starting on port %s\n", port)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("HTTP server failed: %v", err)
	}
}

func statusAlive(w http.ResponseWriter, r *http.Request) {
	if _, err := getCounter(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "service is not running: %s\n", err)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Printf("Alive\n")
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
	counter, _ := getCounter()
	// print
	fmt.Print(line)
	fmt.Printf("Ping / Pongs: %s\n", counter)

	// html out
	fmt.Fprintln(w, line)
	fmt.Fprintf(w, "Ping / Pongs: %s\n", counter)
	printConfigValues(w)

	msg := getGreeterMessage()
	fmt.Fprintf(w, "greetings: %s\n", msg)
}

func getGreeterMessage() string {
	svc := os.Getenv("GREETER_SERVICE")
	if svc == "" {
		return "GREETER_SERVICE not set"
	}
	resp, err := http.Get(fmt.Sprintf("http://%s/", svc))
	if err != nil {
		return fmt.Sprintf("error making request: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Sprintf("received non-OK status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintf("error reading response body: %v", err)
	}

	return string(body)
}
