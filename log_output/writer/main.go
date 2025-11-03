package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
)

var randomUUID string

const iso8601Format = "2006-01-02T15:04:05.000Z"

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

func main() {
	randomUUID = uuid.New().String()
	fmt.Printf("Startup: Generated and stored UUID: %s\n", randomUUID)

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	fmt.Println("Application running. Outputting the UUID every 5 seconds and serving / endpoint (Press Ctrl+C to stop)...")

	for {
		<-ticker.C

		nowUTC := time.Now().UTC()

		timestamp := nowUTC.Format(iso8601Format)

		line := fmt.Sprintf("%s: %s\n", timestamp, randomUUID)

		fmt.Print(line)

		write("/tmp/output.txt", line)
	}
}
