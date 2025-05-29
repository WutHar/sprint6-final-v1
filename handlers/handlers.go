package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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

	originalFilename := header.Filename
	fileExtension := filepath.Ext(originalFilename)

	baseName := fmt.Sprintf("converted_%s", time.Now().UTC().Format("2006-01-02_15-04-05.000000000"))
	fileName := baseName
	if fileExtension == "" {
		fileName = fileName + ".txt"
	} else {
		fileName = fileName + fileExtension
	}

	// Ensure no problematic characters remain in filename from timestamp format
	fileName = strings.ReplaceAll(fileName, ":", "_")
	fileName = strings.ReplaceAll(fileName, "-", "_")
	// If the extension was added, ensure only the last part is the extension
	lastDotIndex := strings.LastIndex(fileName, ".")
	if lastDotIndex != -1 && lastDotIndex < len(fileName)-1 { // Check if a dot exists and is not the very last char
		fileName = strings.ReplaceAll(fileName[:lastDotIndex], ".", "_") + fileName[lastDotIndex:]
	} else {
		fileName = strings.ReplaceAll(fileName, ".", "_") // If no extension, just replace all dots
	}

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
