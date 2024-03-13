package main

import (
	"bytes"
	"compress/gzip"
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

var blockDirPath = ""

func init() {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting working directory:", err)
		return
	}

	blockDirPath = wd + BlockDir
}

func main() {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/api/blocks", getBlocks)
	r.Get("/api/blocks/{id}", findBlockByID)
	r.Put("/api/blocks/{id}", updateBlock)
	r.Post("/api/blocks/", saveBlock)

	if err := http.ListenAndServe(":8081", r); err != nil {
		fmt.Println("Error:", err)
	}
}

func updateBlock(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	file, err := os.OpenFile(blockDirPath+id+".bin", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		http.Error(w, "Error opening file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	defer file.Close()

	writer, err := gzip.NewWriterLevel(file, gzip.BestCompression)
	if err != nil {
		http.Error(w, "Error creating gzip writer: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer writer.Close()

	_, err = writer.Write(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, id)
}

func getBlocks(w http.ResponseWriter, r *http.Request) {
	matches, _ := filepath.Glob(blockDirPath + "*.bin")
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

	// Check if the directory exists, create it if not
	if err := os.MkdirAll(blockDirPath, 0755); err != nil {
		http.Error(w, "Error creating directory: "+err.Error(), http.StatusInternalServerError)
		return
	}

	matches, _ := filepath.Glob(blockDirPath + "g01i*.bin")

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
	// Create the file for writing
	// format is g00i0000.bin
	file, err := os.Create(blockDirPath + filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer file.Close()

	writer, err := gzip.NewWriterLevel(file, gzip.BestCompression)
	defer writer.Close()

	_, err = writer.Write(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, filename)
}

func findBlockByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	file, err := os.Open(blockDirPath + id + ".bin")
	if err != nil {
		http.Error(w, "Error opening file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	defer file.Close()

	reader, err := gzip.NewReader(file)
	if err != nil {
		http.Error(w, "Error creating gzip reader: "+err.Error(), http.StatusInternalServerError)
		return
	}

	reader.Close()

	decompressedData, err := io.ReadAll(reader)
	if err != nil {
		http.Error(w, "Error reading decompressed data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(decompressedData)
}
