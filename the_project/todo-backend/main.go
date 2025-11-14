package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Post struct {
	Body string `json:"body"`
}

var posts = []Post{
	{"Learn JavaScript"},
	{"Learn React"},
	{"Build a project"},
}

func postsGet(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	posts, err := json.Marshal(posts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(posts); err != nil {
		log.Printf("Error writing response: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func postsPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	body := r.FormValue("body")
	post := Post{Body: body}
	posts = append(posts, post)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Post %s was added successfully", body)
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		postsGet(w, r)
		return
	}

	if r.Method == http.MethodPost {
		postsPost(w, r)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		panic("Environmental variable PORT is not set")
	}

	addr := ":" + port

	fmt.Printf("Server started in port %s\n", port)

	http.HandleFunc("/posts", handler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/posts", http.StatusMovedPermanently)
	})

	log.Fatal(http.ListenAndServe(addr, nil))
}
