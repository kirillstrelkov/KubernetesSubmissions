package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

type MyHandler struct {
	Db   *sql.DB
	dbMu sync.RWMutex
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

func (h *MyHandler) handlerAlive(w http.ResponseWriter, r *http.Request) {
	if h.Db == nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "database not connected\n")
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "alive\n")
}

func handleStress(w http.ResponseWriter, _ *http.Request) {
	cores := runtime.NumCPU()
	duration := 60 * time.Second

	fmt.Printf("Starting stress test on %d cores for %v\n", cores, duration)

	done := make(chan bool)

	for range cores {
		go func() {
			for {
				select {
				case <-done:
					return
				default:
					// Spin forever
				}
			}
		}()
	}

	time.Sleep(duration)
	close(done)

	fmt.Fprintf(w, "Finished stress test on %d cores.", cores)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	h := &MyHandler{Db: nil}

	go func() {
		log.Println("Waiting 15 seconds to connect to database...")
		time.Sleep(15 * time.Second)

		db, err := connectToDB()
		if err != nil {
			log.Fatalf("Failed to connect to DB after delay: %v", err)
		}

		err = initDatabase(db)
		if err != nil {
			log.Fatalf("Failed to initialize database after delay: %v", err)
		}

		h.dbMu.Lock()
		h.Db = db
		h.dbMu.Unlock()
		log.Println("Database connected and initialized.")
	}()

	defer func() {
		h.dbMu.RLock()
		dbToClose := h.Db
		h.dbMu.RUnlock()
		if dbToClose != nil {
			if err := dbToClose.Close(); err != nil {
				log.Printf("Error closing database connection: %v", err)
			}
		}
	}()

	addr := ":" + port

	fmt.Printf("Server v2 started in port %s\n", port)

	http.HandleFunc("/pings", h.handlerPings)
	http.HandleFunc("/", h.handler)
	http.HandleFunc("/healthz", h.handlerAlive)
	http.HandleFunc("/stress", handleStress)

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
	h.dbMu.RLock()
	db := h.Db
	h.dbMu.RUnlock()

	if db == nil {
		return fmt.Errorf("database not connected")
	}
	_, err := db.Exec("UPDATE counters SET count = count + 1 WHERE id = 1")
	return err
}

func (h *MyHandler) getCounter() (int, error) {
	h.dbMu.RLock()
	db := h.Db
	h.dbMu.RUnlock()

	if db == nil {
		return 0, fmt.Errorf("database not connected")
	}
	var count int
	row := db.QueryRow("SELECT count FROM counters WHERE id = 1")
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
