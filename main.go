package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
    fs := http.FileServer(http.Dir("./static"))
    http.Handle("/", fs)
    http.HandleFunc("/save-message", func(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
	    http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	    return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
	    http.Error(w, "Failed to read request body", http.StatusInternalServerError)
	    return
	}
	err = os.WriteFile("msg.svg", body, 0644)
	if err != nil {
	    http.Error(w, "Failed to save SVG", http.StatusInternalServerError)
	    return
	}
	fmt.Fprintln(w, "SVG saved successfully")
    })


    http.HandleFunc("/get-svg", func(w http.ResponseWriter, r *http.Request) { 
	if r.Method != http.MethodGet {
	    http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	    return
	}
	svg, err := os.ReadFile("msg.svg")
	if err != nil {
	    http.Error(w, "Failed to read SVG file", http.StatusInternalServerError)
	    return
	}
	w.Header().Set("Content-Type", "image/svg+xml")
	w.Write(svg)
    })

    log.Println("Serving on port 9999")
    err := http.ListenAndServe(":9999", nil)
    if err != nil {
	panic(err)
    }
}
