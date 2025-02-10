package main

import (
	"context"
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Service struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	IconURL string `json:"iconUrl"`
}

type Group struct {
	Name     string    `json:"name"`
	Services []Service `json:"services"`
}

type Dashboard struct {
	Groups []Group `json:"groups"`
}


func loadConfig() (Dashboard, error) {
	var dashboard Dashboard
	configFile, err := os.Open("config/config.json")
	if err != nil {
		return dashboard, err
	}
	defer configFile.Close()

	byteValue, err := io.ReadAll(configFile)
	if err != nil {
		return dashboard, err
	}

	err = json.Unmarshal(byteValue, &dashboard)
	return dashboard, err
}

// handleError is a utility function to handle HTTP error responses
func handleError(w http.ResponseWriter, err error, message string) {
	log.Printf("Error: %v", err)
	http.Error(w, message, http.StatusInternalServerError)
}

// requestHandler handles the HTTP requests and serves the dashboard
func requestHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		handleError(w, err, "Error loading template")
		return
	}

	dashboard, err := loadConfig()
	if err != nil {
		handleError(w, err, "Error loading config")
		return
	}

	if err := tmpl.Execute(w, dashboard); err != nil {
		handleError(w, err, "Error executing template")
	}
}

// startServer initializes and starts the HTTP server
func startServer() *http.Server {
	// Serve static files from the "static" directory
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Serve the dashboard
	http.HandleFunc("/", requestHandler)

	// Start the server
	server := &http.Server{
		Addr:    ":8080",
		Handler: nil, // Default handler, which includes the ones above
	}

	go func() {
		log.Println("Starting server on port 8080")
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	return server
}

// gracefulShutdown waits for a termination signal and shuts down the server gracefully
func gracefulShutdown(server *http.Server) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive a signal
	sigReceived := <-sigChan
	log.Printf("Received signal %v, shutting down server...", sigReceived)

	// Graceful shutdown with a 5-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown failed: %v", err)
	} else {
		log.Println("Server gracefully stopped")
	}
}

func main() {
	// Start the server
	server := startServer()

	// Gracefully handle shutdown
	gracefulShutdown(server)
}
