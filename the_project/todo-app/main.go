package main

import (
	"encoding/json"
	"fmt"
	"html/template"
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
var urlPosts = "http://localhost:8080/posts"

type Post struct {
	Body string `json:"body"`
}

type TemplateData struct {
	Posts []Post
}

func getPosts() []Post {
	var posts []Post

	fmt.Printf("Fetching... %s\n", urlPosts)

	resp, err := http.Get(urlPosts)
	if err != nil {
		fmt.Printf("Error making request: %v", err)
		return posts
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Received non-OK status code: %d", resp.StatusCode)
		return posts

	}

	err = json.NewDecoder(resp.Body).Decode(&posts)
	if err != nil {
		fmt.Printf("Error reading response body: %v", err)
		return posts
	}

	fmt.Printf("Received %d posts\n", len(posts))

	return posts
}

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

	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	posts := getPosts()

	data := TemplateData{
		Posts: posts,
	}

	err := tmpl.Execute(w, data)

	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func main() {
	if url := os.Getenv("POSTS_URL"); url != "" {
		urlPosts = url
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	go fetchInBackground()
	imagePath := fmt.Sprintf("%s/%s", filesRoot, imageFileName)

	if _, err := os.Stat(imagePath); err != nil {
		fetchAndCacheImage()
	} else {
		fmt.Printf("Using cached image: %s\n", imagePath)
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
