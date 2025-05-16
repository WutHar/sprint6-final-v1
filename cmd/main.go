package main

import (
	"log"
	"os"

	"github.com/WutHar/sprint6-final-v1/internal/server"
)

func main() {

	logger := log.New(os.Stdout, "http-server: ", log.LstdFlags)

	srv := server.NewServer(logger)

	if err := srv.Start(); err != nil {

		logger.Fatalf("Ошибка при запуске сервера: %v", err)

	}

}
