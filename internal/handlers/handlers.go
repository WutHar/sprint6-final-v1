package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"github.com/WutHar/sprint6-final-v1/internal/service"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	content, err := os.ReadFile("index.html")
	if err != nil {
		http.Error(w, "Could not read index.html", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(content)
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 * 1024 * 1024)
	if err != nil {
		http.Error(w, "Could not parse form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("myFile")
	if err != nil {
		http.Error(w, "Could not retrieve file from form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Could not read file", http.StatusInternalServerError)
		return
	}
	fileContent := string(fileBytes)
	log.Printf("Содержимое загруженного файла: %q", fileContent)

	convertedText, err := service.DetectAndConvert(fileContent)
	if err != nil {
		http.Error(w, "Could not convert data", http.StatusInternalServerError)
		return
	}

	responseText := convertedText
	if !isMorse(fileContent) {
		responseText = fileContent
	}

	timestamp := time.Now().UTC().String()
	ext := filepath.Ext(header.Filename)
	outputFilename := fmt.Sprintf("converted_%s%s", timestamp, ext)

	err = os.WriteFile(outputFilename, []byte(convertedText), 0644)
	if err != nil {
		log.Printf("Error writing to file: %v", err)
	}

	response := fmt.Sprintf("Конвертированный текст:\n%s\n\nФайл сохранен как: %s", responseText, outputFilename)
	log.Printf("Ответ сервера: %q", response)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(response))
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

func isMorse(data string) bool {
	for _, char := range data {
		if !unicode.Is(unicode.Dash, char) && char != '.' && char != ' ' {
			return false
		}
	}
	return strings.ContainsAny(data, ".-")
}
