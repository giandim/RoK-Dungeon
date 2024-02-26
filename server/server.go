package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

const (
	BlockDir = "/game-engine/assets/blocks/"
)

func getBlocks(w http.ResponseWriter, r *http.Request) {
	wd, err := os.Getwd()
	if err != nil {
		http.Error(w, "Error getting working directory: "+err.Error(), http.StatusInternalServerError)
		return
	}

	dirPath := wd + BlockDir

	matches, _ := filepath.Glob(dirPath + "*.bin")
	var blocks []string

	for _, match := range matches {
		_, filename := filepath.Split(match)
		blocks = append(blocks, strings.Split(filename, ".")[0])
	}

	data, _ := json.Marshal(blocks)

	fmt.Fprint(w, bytes.NewBuffer(data))
}

func saveBlock(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	wd, err := os.Getwd()
	if err != nil {
		http.Error(w, "Error getting working directory: "+err.Error(), http.StatusInternalServerError)
		return
	}

	dirPath := wd + BlockDir

	// Check if the directory exists, create it if not
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		http.Error(w, "Error creating directory: "+err.Error(), http.StatusInternalServerError)
		return
	}

	matches, _ := filepath.Glob(dirPath + "g01i*.bin")

	// Find the most recent block for that group
	maxId := 0
	for _, match := range matches {
		_, filename := filepath.Split(match)

		id, err := strconv.Atoi(strings.Split(strings.Split(filename, "i")[1], ".")[0])
		if err != nil {
			http.Error(w, "Error converting id to integer: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if id > maxId {
			maxId = id
		}
	}

	filename := fmt.Sprintf("g01i%04d.bin", maxId+1)
	// Open the file for writing
	// format is g00i0000.bin
	file, err := os.Create(dirPath + filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, filename)
}

func findBlockByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	wd, err := os.Getwd()
	if err != nil {
		http.Error(w, "Error getting working directory: "+err.Error(), http.StatusInternalServerError)
		return
	}

	file, err := os.ReadFile(wd + BlockDir + id + ".bin")
	if err != nil {

		http.Error(w, "Error opening file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/msgpack")

	w.Write(file)
}

func main() {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Get("/api/blocks", getBlocks)
	r.Get("/api/blocks/{id}", findBlockByID)
	r.Post("/api/blocks", saveBlock)

	if err := http.ListenAndServe(":8081", r); err != nil {
		fmt.Println("Error:", err)
	}
}
