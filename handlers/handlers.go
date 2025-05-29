package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/WutHar/sprint6-final-v1/service"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./index.html")
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse form: %v", err), http.StatusInternalServerError)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not get file from form: %v", err), http.StatusBadRequest)
		return
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not read file content: %v", err), http.StatusInternalServerError)
		return
	}

	convertedContent, err := service.DetectAndConvert(string(content))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error converting content: %v", err), http.StatusInternalServerError)
		return
	}

	fileExtension := filepath.Ext(header.Filename)
	fileName := fmt.Sprintf("converted_%s%s", time.Now().UTC().Format("2006-01-02_15-04-05.000000000"), fileExtension)

	outputFile, err := os.Create(fileName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not create output file: %v", err), http.StatusInternalServerError)
		return
	}
	defer outputFile.Close()

	_, err = outputFile.WriteString(convertedContent)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not write to output file: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(convertedContent))
}
