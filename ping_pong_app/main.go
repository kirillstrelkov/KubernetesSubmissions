package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

var counter int

func handler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		counter++
	}()

	response := fmt.Sprintf("pong %v\n", counter)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, response)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := ":" + port

	fmt.Printf("Server started in port %s\n", port)

	http.HandleFunc("/pingpong", handler)

	log.Fatal(http.ListenAndServe(addr, nil))
}
