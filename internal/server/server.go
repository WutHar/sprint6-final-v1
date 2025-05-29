package server

import (
	"log"
	"net/http"
	"time"

	"github.com/WutHar/sprint6-final-v1/internal/handlers"
)

type Server struct {
	httpServer *http.Server
	logger     *log.Logger
}

func NewServer(logger *log.Logger) *Server {
	router := http.NewServeMux()

	// Регистрируем наши хендлеры.
	router.HandleFunc("/", handlers.IndexHandler)
	router.HandleFunc("/upload", handlers.UploadHandler)

	httpServer := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	return &Server{
		httpServer: httpServer,
		logger:     logger,
	}
}

func (s *Server) Start() error {
	s.logger.Printf("Запуск сервера на http://localhost%s", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}
