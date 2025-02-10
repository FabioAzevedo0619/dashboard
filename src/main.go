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

func requestHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Printf("Error loading template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	dashboard, err := loadConfig()
	if err != nil {
		log.Printf("Error loading config: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, dashboard); err != nil {
		log.Printf("Error executing template: %v", err)
	}
}

func main() {
	// Serve static files from the "static" directory
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Serve the dashboard
	http.HandleFunc("/", requestHandler)

	// Start the server
	server := &http.Server{
		Addr:    ":8080",
		Handler: nil, // Default handler, which includes the ones above
	}

	// Start the server in a separate goroutine
	go func() {
		log.Println("Starting server on port 8080")
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	// Wait for termination signal for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive a signal
	<-sigChan
	log.Println("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server Shutdown Failed:%v", err)
	} else {
		log.Println("Server gracefully stopped")
	}
}
