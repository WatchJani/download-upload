package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	http.HandleFunc("GET /download", Download)
	http.HandleFunc("POST /upload", Upload)

	http.ListenAndServe(":5000", nil)
}

func Upload(w http.ResponseWriter, r *http.Request) {
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	dst, err := os.Create(filepath.Join(".", handler.Filename))
	if err != nil {
		http.Error(w, "Error creating file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}

	// Send success response
	fmt.Fprintf(w, "File uploaded successfully: %s\n", handler.Filename)
}

func Download(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("add.svg")
	if os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		http.Error(w, "Error retrieving file info", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileInfo.Name()))
	http.ServeContent(w, r, fileInfo.Name(), fileInfo.ModTime(), file)
}
