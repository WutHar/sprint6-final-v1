package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/WutHar/sprint6-final/service"
)

// IndexHandler обрабатывает запрос к корневому пути и возвращает содержимое index.html.
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
	_, err = w.Write(content)
	if err != nil {
		// Логирование ошибки записи в ResponseWriter
		fmt.Println("Error writing response:", err)
	}
}

// UploadHandler обрабатывает загрузку файла, конвертирует его содержимое и сохраняет в новый файл.
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

	// Получаем файл из формы. "file" - это имя поля в вашей HTML-форме.
	file, header, err := r.FormFile("myFile")
	if err != nil {
		http.Error(w, "Could not retrieve file from form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Читаем содержимое файла.
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Could not read file", http.StatusInternalServerError)
		return
	}
	fileContent := string(fileBytes)

	// Определяем тип и конвертируем содержимое.
	convertedText, err := service.DetectAndConvert(fileContent)
	if err != nil {
		http.Error(w, "Could not detect and convert content", http.StatusInternalServerError)
		return
	}

	// Создаем имя для нового файла.
	timestamp := time.Now().UTC().String()
	ext := filepath.Ext(header.Filename)
	newFilename := fmt.Sprintf("converted_%s%s", timestamp, ext)

	// Создаем и записываем сконвертированный текст в новый файл.
	newFile, err := os.Create(newFilename)
	if err != nil {
		http.Error(w, "Could not create output file", http.StatusInternalServerError)
		return
	}
	defer newFile.Close()

	_, err = newFile.WriteString(convertedText)
	if err != nil {
		http.Error(w, "Could not write to output file", http.StatusInternalServerError)
		return
	}

	// Возвращаем сконвертированный текст в ответе.
	w.WriteHeader(http.StatusOK)
	_, err = fmt.Fprintf(w, "Конвертированный текст:\n%s\nФайл сохранен как: %s\n", convertedText, newFilename)
	if err != nil {
		// Логирование ошибки записи в ResponseWriter
		fmt.Println("Error writing response:", err)
	}
}
