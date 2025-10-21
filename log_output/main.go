package main

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

func main() {
	randomUUID := uuid.New().String()
	fmt.Printf("Startup: Generated and stored UUID: %s\n", randomUUID)

	const iso8601Format = "2006-01-02T15:04:05.000Z"

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	fmt.Println("Application running. Outputting the UUID every 5 seconds (Press Ctrl+C to stop)...")

	for {
		<-ticker.C

		nowUTC := time.Now().UTC()

		timestamp := nowUTC.Format(iso8601Format)

		fmt.Printf("%s: %s\n", timestamp, randomUUID)
	}
}
