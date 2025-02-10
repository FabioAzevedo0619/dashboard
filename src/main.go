package main

import (
	"fmt"
	"html/template"
	"net/http"
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

func helloHandler(w http.ResponseWriter, r *http.Request) {
    tmpl := template.Must(template.ParseFiles("template.html"))
    data := struct {
        Title   string
        Message string
    }{
        Title:   "Hello Page",
        Message: "Hello, World!",
    }
    tmpl.Execute(w, data)
}

func main() {
    http.HandleFunc("/", helloHandler)
    fmt.Println("Starting server at port 8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        fmt.Println("Error starting server:", err)
    }
}