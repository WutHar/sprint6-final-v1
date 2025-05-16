package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/WutHar/sprint6-final-v1/internal/service"
)

// IndexHandler обрабатывает запросы к корневому пути и возвращает HTML-форму.
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

// UploadHandler обрабатывает загрузку файла, определяет его тип и конвертирует.
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Парсим multipart form, с ограничением в 10 MB.
	err := r.ParseMultipartForm(10 * 1024 * 1024)
	if err != nil {
		http.Error(w, "Could not parse form", http.StatusBadRequest)
		return
	}

	// Получаем файл из формы
	file, header, err := r.FormFile("myFile")
	if err != nil {
		http.Error(w, "Could not retrieve file from form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Читаем данные из файла
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Could not read file", http.StatusInternalServerError)
		return
	}
	fileContent := string(fileBytes)

	// Определяем тип и конвертируем данные
	convertedText, err := service.DetectAndConvert(fileContent)
	if err != nil {
		http.Error(w, "Could not convert data", http.StatusInternalServerError)
		return
	}

	// Создаем имя для локального файла
	timestamp := time.Now().UTC().String()
	ext := filepath.Ext(header.Filename)
	outputFilename := fmt.Sprintf("converted_%s%s", timestamp, ext)

	// Записываем результат в локальный файл
	err = os.WriteFile(outputFilename, []byte(convertedText), 0644)
	if err != nil {
		log.Printf("Error writing to file: %v", err)
	}

	// Возвращаем результат конвертации
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(fmt.Sprintf("Конвертированный текст:\n%s\n\nФайл сохранен как: %s", convertedText, outputFilename)))
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}
