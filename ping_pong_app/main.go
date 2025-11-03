package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

var counter int

func write(filePath, line string) {
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if _, err := f.WriteString(line); err != nil {
		log.Fatal(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		counter++
	}()

	write("/tmp/counter.txt", fmt.Sprintf("%d", counter))

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
