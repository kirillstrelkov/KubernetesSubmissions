package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

const (
	filesRoot     = "/tmp/image"
	imageFileName = "image.jpg"
)

var isOldImage = false

func todoHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	html := `<!DOCTYPE html>
<html>
<head>
    <title>The project App</title>
</head>
<body>
    <h1>The project App</h1>
	</br>
	<img src="/image" alt="Random image" width="200" height="200"/>
	</br>
	<p>DevOps with Kubernetes 2025</p>
</body>
</html>`
	fmt.Fprint(w, html)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	go fetchInBackground()
	imagePath := fmt.Sprintf("%s/%s", filesRoot, imageFileName)

	if _, err := os.Stat(imagePath); err != nil {
		fetchAndCacheImage()
	}

	addr := ":" + port

	fmt.Printf("Server started in port %s\n", port)

	http.HandleFunc("/image", imageHandler)
	http.HandleFunc("/", todoHandler)

	log.Fatal(http.ListenAndServe(addr, nil))

}

func fetchInBackground() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for {
		<-ticker.C

		isOldImage = true
	}

}

func imageHandler(w http.ResponseWriter, r *http.Request) {
	if isOldImage {
		go fetchAndCacheImage()
	}

	imagePath := fmt.Sprintf("%s/%s", filesRoot, imageFileName)

	http.ServeFile(w, r, imagePath)
}

func fetchAndCacheImage() {
	log.Println("Fetching new image from picsum.photos...")

	id := rand.Intn(1000) + 1
	url := fmt.Sprintf("https://picsum.photos/%d", id)
	log.Printf("Image URL: %s", url)

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to fetch image: %v", err)
		return
	}
	defer resp.Body.Close()

	err = os.MkdirAll(filesRoot, os.ModePerm)
	if err != nil {
		log.Printf("Failed to create directory: %v", err)
		return
	}

	imagePath := fmt.Sprintf("%s/%s", filesRoot, imageFileName)
	out, err := os.Create(imagePath)
	if err != nil {
		log.Printf("Failed to create image file: %v", err)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Printf("Failed to save image: %v", err)
	}
	log.Println("Successfully fetched and cached new image.")

	isOldImage = false
}
