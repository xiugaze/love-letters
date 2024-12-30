package main

import (
	"log"
	"net/http"
)

func main() {
    fs := http.FileServer(http.Dir("./static"))
    http.Handle("/", fs)

    err := http.ListenAndServe(":9999", nil)
    if err != nil {
        panic(err)
    }
    log.Println("Serving on port 9999")
}
