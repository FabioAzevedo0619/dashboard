package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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

    byteValue, err := ioutil.ReadAll(configFile)
    if err != nil {
        return dashboard, err
    }

    err = json.Unmarshal(byteValue, &dashboard)
    return dashboard, err
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
    tmpl := template.Must(template.ParseFiles("templates/index.html"))

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
    // Serve static files from the "src/static" directory
    fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

    http.HandleFunc("/", helloHandler)
    fmt.Println("Starting server at port 8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        fmt.Println("Error starting server:", err)
    }
}