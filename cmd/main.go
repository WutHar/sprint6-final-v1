package main

import (
	"log"
	"os"

	"github.com/WutHar/sprint6-final/server"
)

func main() {
	// Создаем логгер.
	logger := log.New(os.Stdout, "http-server: ", log.LstdFlags)

	// Создаем экземпляр сервера.
	srv := server.NewServer(logger)

	// Запускаем сервер.
	if err := srv.Start(); err != nil {
		logger.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
