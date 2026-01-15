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
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		msg := os.Getenv("MSG")
		if msg == "" {
			msg = "Hello"
		}
		fmt.Fprintf(w, "%s", msg)
	})
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
