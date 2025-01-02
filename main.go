package main

import (
  "fmt"
  "io"
  "log"
  "net/http"
  "os"
  "time"
)

const DEFAULT_PORT string = "8080"
var static string
var media string
var port string

func get_vars() (string, string, string) {
  static := os.Getenv("STATIC_FILES_PATH")
  if static == "" {
    static = "./static"
  }
  media := os.Getenv("MEDIA_FILES_PATH")
  if media == "" {
    media = "./media"
  }

  port := os.Getenv("PORT")
  if port == "" {
    port = DEFAULT_PORT
  }
  return static, media, port
}

/* handle post request containing svg */
func save_handler(w http.ResponseWriter, r *http.Request) {
  // go uses 2006-01-02 15:04:05 for time format referencing
  time := time.Now().Format("2006-01-02")
  if r.Method != http.MethodPost {
    http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
    return
  }

  body, err := io.ReadAll(r.Body)
  if err != nil {
    http.Error(w, "Failed to read request body (invalid .svg?)", http.StatusInternalServerError)
  }
  path := media + "/" + time + ".svg";
  log.Println("Attempting write to " + path)
  err = os.WriteFile(path, body, 0644) // rw-r--r--
  if err != nil {
    http.Error(w, "Failed to save SVG to " + path, http.StatusInternalServerError)
    fmt.Println(err)
    return
  }
  fmt.Fprintln(w, "saved successfully")
}

func get_handler(w http.ResponseWriter, r *http.Request) {
  time := time.Now().Format("2006-01-02")
  if r.Method != http.MethodGet {
    http.Error(w, "Invalid Request method", http.StatusMethodNotAllowed)
    return
  }
  svg, err := os.ReadFile(media + "/" + time + ".svg")
  if err != nil {
    http.Error(w, "Failed to read SVG file (does file exist?)", http.StatusInternalServerError)
    return
  }
  w.Header().Set("Content-Type", "image/svg+xml")
  w.Write(svg)
}

func main() {
  static, media, port = get_vars()
  http.Handle("/", http.FileServer(http.Dir(static)))

  http.HandleFunc("/save", save_handler)
  http.HandleFunc("/get", get_handler)

  log.Println("Serving on port " + port)
  log.Println("media: " + media)
  err := http.ListenAndServe(":" + port, nil)
  if err != nil {
    panic(err)
  }
}
