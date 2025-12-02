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
	"path/filepath"
	"time"
)

var isOldImage = false
var alive = false

type Post struct {
	Body string `json:"body"`
}

type TemplateData struct {
	Posts []Post
}

type EnvVars struct {
	Port     string
	PostsURL string
	ImgPath  string
	ImgURL   string
}

type MyHandler struct {
	Vars EnvVars
}

func NewHandler() *MyHandler {
	var missing_vars []string

	port := os.Getenv("PORT")
	if port == "" {
		missing_vars = append(missing_vars, "PORT")
	}

	urlPosts := os.Getenv("POSTS_URL")
	if urlPosts == "" {
		missing_vars = append(missing_vars, "POSTS_URL")
	}

	imgPath := os.Getenv("IMG_PATH")
	if imgPath == "" {
		missing_vars = append(missing_vars, "IMG_PATH")
	}

	imgURL := os.Getenv("IMG_URL")
	if imgURL == "" {
		missing_vars = append(missing_vars, "IMG_URL")
	}

	if len(missing_vars) > 0 {
		panic(fmt.Sprintf("Environment variables %v is not set", &missing_vars))
	}

	vars := EnvVars{
		Port:     port,
		PostsURL: urlPosts,
		ImgPath:  imgPath,
		ImgURL:   imgURL,
	}
	return &MyHandler{Vars: vars}
}

func getPosts(urlPosts string) []Post {
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

func createImageFile(h *MyHandler) {
	fmt.Println("Waiting 5 seconds before creating image file...")
	time.Sleep(5 * time.Second)

	imgPath := h.Vars.ImgPath
	if _, err := os.Stat(imgPath); err != nil {
		h.fetchAndCacheImage()
	} else {
		fmt.Printf("Using cached image: %s\n", imgPath)
	}
	alive = true
}

func handleAlive(w http.ResponseWriter, r *http.Request) {
	if !alive {
		http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "alive")
}

func (h *MyHandler) todoHandler(w http.ResponseWriter, r *http.Request) {
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
	posts := getPosts(h.Vars.PostsURL)

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
	h := NewHandler()

	go createImageFile(h)

	http.HandleFunc("/image", h.imageHandler)
	http.HandleFunc("/", h.todoHandler)
	http.HandleFunc("/healthz", handleAlive)

	addr := ":" + h.Vars.Port
	fmt.Printf("Server started in port %s\n", h.Vars.Port)

	go fetchInBackground()

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

func (h *MyHandler) imageHandler(w http.ResponseWriter, r *http.Request) {
	if isOldImage {
		go h.fetchAndCacheImage()
	}

	http.ServeFile(w, r, h.Vars.ImgPath)
}

func (h *MyHandler) fetchAndCacheImage() {
	log.Println("Fetching new image...")

	id := rand.Intn(1000) + 1
	url := fmt.Sprintf(h.Vars.ImgURL, id)
	log.Printf("Image URL: %s", url)

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to fetch image: %v", err)
		return
	}
	defer resp.Body.Close()

	imagePath := h.Vars.ImgPath
	folder := filepath.Dir(imagePath)
	err = os.MkdirAll(folder, os.ModePerm)
	if err != nil {
		log.Printf("Failed to create directory: %v", err)
		return
	}

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
