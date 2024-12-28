package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Priyang1310/Students-API-GO/internal/config"
	"github.com/Priyang1310/Students-API-GO/internal/http/handlers/student"
	"github.com/Priyang1310/Students-API-GO/internal/storage/sqlite"
)

func main() {
	//load config
	cfg := config.MustLoad()

	//database setup
	storage, err := sqlite.New(cfg)

	if err != nil {
		log.Fatal(err)
	}

	slog.Info("Storage Initialized!", slog.String("env", cfg.Env))

	//setup router
	router := http.NewServeMux()

	router.HandleFunc("POST /api/students", student.New(storage))

	// Setup the HTTP server with the specified address and handler
	server := http.Server{
		Addr:    cfg.Addr, // The address the server will listen on (e.g., ":3000")
		Handler: router,   // The router that will handle incoming requests
	}

	// Create a channel to listen for OS signals (like interrupt or termination)
	done := make(chan os.Signal, 1)

	// Notify the channel when an interrupt (Ctrl+C) or termination signal is received
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Start the server in a separate goroutine
	go func() {
		// Listen and serve HTTP requests
		errr := server.ListenAndServe()

		// If there's an error starting the server, log it and exit
		if errr != nil {
			log.Fatalln("Error starting server:", errr.Error())
		}
	}()

	// Wait for a signal to be received
	<-done
	fmt.Println("\n=======================")
	log.Println("Server gracefully stopped")
}
