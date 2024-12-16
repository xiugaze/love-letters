package main

import (
	"log"
	"net/http"
)

func main() {
    fs := http.FileServer(http.Dir("./static"))
    http.Handle("/", fs)

    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        panic(err)
    }
    log.Println("Serving on port 8080")
}
