package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func todoHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	greeting := os.Getenv("GREETING")
	if greeting == "" {
		greeting = "Hello World"
	}

	html := `<!DOCTYPE html>
<html>
<head>
    <title>Todo App</title>
</head>
<body>
    <h1>%s</h1>
</body>
</html>`
	fmt.Fprintf(w, html, greeting)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := ":" + port

	fmt.Printf("Server started in port %s\n", port)

	http.HandleFunc("/", todoHandler)

	log.Fatal(http.ListenAndServe(addr, nil))
}
