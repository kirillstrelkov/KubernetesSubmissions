package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
)

type Post struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
	Done bool   `json:"done"`
}

type MyHandler struct {
	Db   *sql.DB
	Nats *nats.Conn
	dbMu sync.RWMutex
}

func connectToNATS() *nats.Conn {
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		panic("Environmental variable NATS_URL is not set")
	}

	log.Printf("Connecting to NATS at %s...", natsURL)
	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Printf("WARNING: Failed to connect to NATS: %v. The app will work, but messages won't be sent.", err)
		return nil
	}

	log.Println("Successfully connected to NATS.")
	return nc
}

func publishToNATS(nc *nats.Conn, subject string, payload any) {
	msgBytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling NATS payload for subject %s: %v", subject, err)
		return
	}

	if err := nc.Publish(subject, msgBytes); err != nil {
		log.Printf("Error publishing to NATS subject %s: %v", subject, err)
	} else {
		log.Printf("NATS: Published to '%s': %s", subject, string(msgBytes))
	}
}

func getPosts(db *sql.DB) ([]Post, error) {
	rows, err := db.Query("SELECT id, body, done FROM posts ORDER BY id ASC")
	if err != nil {
		return nil, fmt.Errorf("error querying posts: %w", err)
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.Body, &post.Done); err != nil {
			return nil, fmt.Errorf("error scanning post row: %w", err)
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over post rows: %w", err)
	}

	return posts, nil
}

func (h *MyHandler) getDB() *sql.DB {
	h.dbMu.RLock()
	defer h.dbMu.RUnlock()
	return h.Db
}

func (h *MyHandler) postsGet(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	db := h.getDB()

	if db == nil {
		http.Error(w, "Database not ready", http.StatusServiceUnavailable)
		return
	}

	posts, err := getPosts(db)
	if err != nil {
		log.Printf("Error retrieving posts: %v", err)
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

	db := h.getDB()

	if db == nil {
		http.Error(w, "Database not ready", http.StatusServiceUnavailable)
		return
	}

	log.Printf("Adding a new todo: %s", body)

	var newID int
	err := db.QueryRow("INSERT INTO posts (body) VALUES ($1) RETURNING id", body).Scan(&newID)
	if err != nil {
		log.Printf("Error inserting post into database: %v", err)
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	newPost := Post{ID: newID, Body: body, Done: false}
	publishToNATS(h.Nats, "todo.created", newPost)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Post %s was added successfully", body)
}

func (h *MyHandler) handleAlive(w http.ResponseWriter, r *http.Request) {
	if h.Db == nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "database not connected\n")
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "alive")
}

func (h *MyHandler) markDoneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	db := h.getDB()

	if db == nil {
		http.Error(w, "Database not ready", http.StatusServiceUnavailable)
		return
	}

	result, err := db.Exec("UPDATE posts SET done = TRUE WHERE id = $1", id)
	if err != nil {
		log.Printf("Error updating post in database: %v", err)
		http.Error(w, "Failed to update post", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected: %v", err)
		http.Error(w, "Failed to update post", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	updatedPost := Post{ID: id, Done: true}
	publishToNATS(h.Nats, "todo.updated", updatedPost)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Post with ID %d marked as done", id)
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

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		panic("Environmental variable PORT is not set")
	}

	nc := connectToNATS()
	if nc != nil {
		defer nc.Close()
	}

	h := &MyHandler{
		Db:   nil,
		Nats: nc,
	}

	go func() {
		log.Println("Waiting 5 seconds to connect to database...")
		time.Sleep(5 * time.Second)

		db, err := connectToDB()
		if err != nil {
			log.Fatalf("Failed to connect to DB: %v", err)
		}

		err = initDatabase(db)
		if err != nil {
			log.Fatalf("Failed to initialize database: %v", err)
		}

		h.dbMu.Lock()
		h.Db = db
		h.dbMu.Unlock()
		log.Println("Database connected and initialized.")
	}()

	defer func() {
		dbToClose := h.getDB()

		if dbToClose != nil {
			if err := dbToClose.Close(); err != nil {
				log.Printf("Error closing database connection: %v", err)
			}
		}
	}()

	addr := ":" + port

	fmt.Printf("Server v2 started in port %s\n", port)

	mux := http.NewServeMux()

	mux.HandleFunc("/posts", h.handler)
	mux.HandleFunc("/healthz", h.handleAlive)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/posts", http.StatusMovedPermanently)
	})
	mux.HandleFunc("/todos/{id}", h.markDoneHandler)

	log.Fatal(http.ListenAndServe(addr, enableCORS(mux)))
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
			body TEXT NOT NULL,
			done BOOLEAN DEFAULT FALSE
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create posts table: %w", err)
	}

	_, err = db.Exec(`ALTER TABLE posts ADD COLUMN IF NOT EXISTS done BOOLEAN DEFAULT FALSE;`)
	if err != nil {
		return fmt.Errorf("failed to add done column: %w", err)
	}

	posts, _ := getPosts(db)
	if len(posts) == 0 {
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
	}

	log.Println("Posts table initialized successfully.")
	return nil
}
