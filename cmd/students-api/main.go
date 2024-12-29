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

// The main function is the entry point of the program.
func main() {
	// Load the configuration from the environment or a configuration file.
	// The config.MustLoad function returns a Config object or panics if there's an error.
	cfg := config.MustLoad()

	// Initialize the database storage using the provided configuration.
	// The sqlite.New function returns a Storage object or an error if the database cannot be initialized.
	storage, err := sqlite.New(cfg)

	// If there's an error initializing the database, log the error and exit the program.
	if err != nil {
		log.Fatal(err)
	}

	// Log a message indicating that the storage has been initialized.
	slog.Info("Storage Initialized!", slog.String("env", cfg.Env))

	// Create a new HTTP request multiplexer to handle incoming requests.
	router := http.NewServeMux()

	// Define the routes for the API endpoints.
	// Each route is associated with a specific handler function that will be called when the route is accessed.
	router.HandleFunc("POST /api/students", student.New(storage))               // Create a new student
	router.HandleFunc("GET /api/students/{id}", student.GetById(storage))       // Get a student by ID
	router.HandleFunc("GET /api/students", student.GetAll(storage))             // Get all students
	router.HandleFunc("PUT /api/students/{id}", student.Update(storage))        // Update a student
	router.HandleFunc("DELETE /api/students/{id}", student.DeleteById(storage)) // Delete a student by ID
	router.HandleFunc("DELETE /api/students", student.DeleteAll(storage))       // Delete all students

	// Create a new HTTP server with the specified address and handler.
	// The server will listen for incoming requests on the specified address and route them to the associated handler functions.
	server := http.Server{
		Addr:    cfg.Addr, // The address the server will listen on (e.g., ":3000")
		Handler: router,   // The router that will handle incoming requests
	}

	// Create a channel to listen for OS signals (like interrupt or termination).
	done := make(chan os.Signal, 1)

	// Notify the channel when an interrupt (Ctrl+C) or termination signal is received.
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Start the server in a separate goroutine.
	// This allows the server to run concurrently with the main goroutine.
	go func() {
		// Listen and serve HTTP requests.
		// The server will continue to run until it is stopped or an error occurs.
		errr := server.ListenAndServe()

		// If there's an error starting the server, log the error and exit the program.
		if errr != nil {
			log.Fatalln("Error starting server:", errr.Error())
		}
	}()

	// Wait for a signal to be received.
	// This will block the main goroutine until a signal is received.
	<-done

	// Log a message indicating that the server has been stopped.
	fmt.Println("\n=======================")
	log.Println("Server gracefully stopped")
}
