package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

type MyHandler struct {
	Db *sql.DB
}

func (h *MyHandler) handler(w http.ResponseWriter, r *http.Request) {
	count, err := h.getCounter()
	if err != nil {
		log.Printf("Error getting counter: %v\n", err)
		http.Error(w, "Failed to get counter", http.StatusInternalServerError)
		return
	}

	response := fmt.Sprintf("pong %v\n", count)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	err = h.incrementCounter()
	if err != nil {
		log.Printf("Error incrementing counter: %v\n", err)
		http.Error(w, "Failed to increment counter", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, response)
}

func (h *MyHandler) handlerPings(w http.ResponseWriter, r *http.Request) {
	count, err := h.getCounter()
	if err != nil {
		log.Printf("Error getting counter: %v\n", err)
		http.Error(w, "Failed to get counter", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%d", count)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	db, err := connectToDB()
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	h := &MyHandler{Db: db}
	defer h.Db.Close()

	err = initDatabase(h.Db)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	addr := ":" + port

	fmt.Printf("Server started in port %s\n", port)

	http.HandleFunc("/pingpong", h.handler)
	http.HandleFunc("/pings", h.handlerPings)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func connectToDB() (*sql.DB, error) {
	connStr := os.Getenv("DB_URL")
	if connStr == "" {
		panic("DB_URL is not set")
	}
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func initDatabase(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS counters (
			id SERIAL PRIMARY KEY,
			count INTEGER DEFAULT 0
		);
		INSERT INTO counters (count) SELECT 0 WHERE NOT EXISTS (SELECT 1 FROM counters);
	`)
	return err
}

func (h *MyHandler) incrementCounter() error {
	_, err := h.Db.Exec("UPDATE counters SET count = count + 1 WHERE id = 1")
	return err
}

func (h *MyHandler) getCounter() (int, error) {
	var count int
	row := h.Db.QueryRow("SELECT count FROM counters WHERE id = 1")
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
