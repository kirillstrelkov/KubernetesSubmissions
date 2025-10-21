package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func todoHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "Welcome to the Todo Application Server!\nPath requested: %s", r.URL.Path)
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
