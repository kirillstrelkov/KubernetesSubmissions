package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/nats-io/nats.go"
)

type TodoIncoming struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
	Done bool   `json:"done"`
}

type TelegramMessage struct {
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
}

var appEnv string

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		panic("Environmental variable PORT is not set")
	}

	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		panic("Environmental variable NATS_URL is not set")
	}

	appEnv = os.Getenv("APP_ENV")
	if appEnv == "" {
		panic("Environmental variable APP_ENV is not set")
	}

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	chatID := os.Getenv("TELEGRAM_CHAT_ID")

	if appEnv == "production" && (botToken == "" || chatID == "") {
		panic("TELEGRAM_BOT_TOKEN and TELEGRAM_CHAT_ID environment variables must be set in production")
	}

	log.Printf("Starting broadcaster in %s mode", appEnv)
	log.Printf("Connecting to NATS at %s...", natsURL)
	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()
	log.Println("Connected to NATS.")

	hostname, _ := os.Hostname()
	user := os.Getenv("USER")
	if user == "" {
		user = "broadcaster"
	}

	log.Println("Subscribing to todo.created and todo.updated subjects...")
	_, err = nc.QueueSubscribe("todo.created", "broadcaster_workers", func(m *nats.Msg) {
		sendMessage(nc, m, "created", botToken, chatID, user, hostname)
	})
	if err != nil {
		log.Fatal(err)
	}

	_, err = nc.QueueSubscribe("todo.updated", "broadcaster_workers", func(m *nats.Msg) {
		sendMessage(nc, m, "updated", botToken, chatID, user, hostname)
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Broadcaster is running. Waiting for messages...")

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "alive\n")
	})
	addr := ":" + port
	fmt.Printf("Server started in port %s\n", port)

	http.ListenAndServe(addr, mux)
}

func sendMessage(nc *nats.Conn, m *nats.Msg, action string, token, chatID, user, hostname string) {
	var todo TodoIncoming

	if err := json.Unmarshal(m.Data, &todo); err != nil {
		log.Printf("Error unmarshalling JSON: %v", err)
		return
	}

	log.Printf("Received todo.%s: %s", action, todo.Body)

	prettyJSON, _ := json.MarshalIndent(todo, "", "  ")

	message := fmt.Sprintf("A todo was %s:\n%s\n\nbroadcasted by %s @ %s",
		action,
		string(prettyJSON),
		user,
		hostname,
	)
	if appEnv == "staging" {
		log.Printf("[STAGING] Skip sending message to Telegram:\n```%s\n```", message)
		return
	}

	sendToTelegram(token, chatID, message)
}

func sendToTelegram(token string, chatID string, text string) {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)

	payload := TelegramMessage{
		ChatID: chatID,
		Text:   text,
	}

	jsonPayload, _ := json.Marshal(payload)

	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		log.Printf("Failed to send to Telegram: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Telegram returned non-OK status: %d", resp.StatusCode)
	}
}
