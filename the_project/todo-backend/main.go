package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

type Post struct {
	Body string `json:"body"`
}

type MyHandler struct {
	Db *sql.DB
}

func (h *MyHandler) postsGet(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := h.Db.Query("SELECT body FROM posts ORDER BY id ASC")
	if err != nil {
		log.Printf("Error querying posts: %v", err)
		http.Error(w, "Failed to retrieve posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.Body); err != nil {
			log.Printf("Error scanning post row: %v", err)
			http.Error(w, "Failed to process posts", http.StatusInternalServerError)
			return
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over post rows: %v", err)
		http.Error(w, "Failed to retrieve posts", http.StatusInternalServerError)
		return
	}

	jsonPosts, err := json.Marshal(posts)
	if err != nil {
		log.Printf("Error marshalling posts to JSON: %v", err)
		http.Error(w, "Failed to serialize posts", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(jsonPosts); err != nil {
		log.Printf("Error writing response: %v", err)
		// No http.Error here as headers might already be sent
	}
}

func (h *MyHandler) postsPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Printf("Error parsing form: %v", err)
		http.Error(w, "Bad request: failed to parse form", http.StatusBadRequest)
		return
	}

	body := r.FormValue("body")
	if body == "" {
		http.Error(w, "Bad request: 'body' cannot be empty", http.StatusBadRequest)
		return
	}

	if len(body) > 140 {
		log.Printf("Todo is too long: %d characters", len(body))
		http.Error(w, "Bad request: 'body' cannot be longer than 140 characters", http.StatusBadRequest)
		return
	}

	log.Printf("Adding a new todo: %s", body)
	_, err := h.Db.Exec("INSERT INTO posts (body) VALUES ($1)", body)
	if err != nil {
		log.Printf("Error inserting post into database: %v", err)
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Post %s was added successfully", body)
}

func (h *MyHandler) handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.postsGet(w, r)
		return
	}

	if r.Method == http.MethodPost {
		h.postsPost(w, r)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		panic("Environmental variable PORT is not set")
	}

	db, err := connectToDB()
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	h := &MyHandler{Db: db}
	defer func() {
		if err := h.Db.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}()

	err = initDatabase(h.Db)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	addr := ":" + port

	fmt.Printf("Server started in port %s\n", port)

	http.HandleFunc("/posts", h.handler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/posts", http.StatusMovedPermanently)
	})

	log.Fatal(http.ListenAndServe(addr, nil))
}

func connectToDB() (*sql.DB, error) {
	connStr := os.Getenv("DB_URL")
	if connStr == "" {
		return nil, fmt.Errorf("DB_URL environment variable is not set")
	}
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	log.Println("Successfully connected to the database.")
	return db, nil
}

func initDatabase(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS posts (
			id SERIAL PRIMARY KEY,
			body TEXT NOT NULL
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create posts table: %w", err)
	}

	var posts = []string{
		"Learn JavaScript",
		"Learn React",
		"Build a project",
	}
	for _, post := range posts {
		_, err := db.Exec("INSERT INTO posts (body) VALUES ($1)", post)
		if err != nil {
			log.Printf("Error inserting post into database: %v", err)
			panic(err)
		}
	}

	log.Println("Posts table initialized successfully.")
	return nil
}
