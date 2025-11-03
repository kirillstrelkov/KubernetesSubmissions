package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
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

func readFile(filePath string) string {
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	return string(data)
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
	response := readFile("/tmp/output.txt")
	if response == "" {
		w.WriteHeader(http.StatusNoContent)
		fmt.Fprint(w, "No content")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Printf("File content: %s", response)
	fmt.Fprint(w, response)
}
