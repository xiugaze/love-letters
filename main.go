package main

import (
	"bytes"
	"encoding/base64"
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

func save_handler(w http.ResponseWriter, r *http.Request) {
    time := time.Now().Format("2006-01-02")
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }
    
    body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Failed to read request body", http.StatusInternalServerError)
        return
    }
    
    bmpData := bytes.Split(body, []byte(","))[1]
    decodedData, err := base64.StdEncoding.DecodeString(string(bmpData))
    if err != nil {
        http.Error(w, "Failed to decode BMP data", http.StatusInternalServerError)
        return
    }
    
    path := media + "/" + time + ".bmp"
    log.Println("Attempting write to " + path)
    err = os.WriteFile(path, decodedData, 0644)
    if err != nil {
        http.Error(w, "Failed to save BMP to "+path, http.StatusInternalServerError)
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
    bmpData, err := os.ReadFile(media + "/" + time + ".bmp")
    if err != nil {
        http.Error(w, "Failed to read BMP file (does file exist?)", http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "image/bmp")
    w.Write(bmpData)
}

func get_bin_handler(w http.ResponseWriter, r *http.Request) {
    time := time.Now().Format("2006-01-02")
    if r.Method != http.MethodGet {
        http.Error(w, "Invalid Request method", http.StatusMethodNotAllowed)
        return
    }

    bmpData, err := os.ReadFile(media + "/" + time + ".bmp")
    if err != nil {
        http.Error(w, "Failed to read BMP file", http.StatusInternalServerError)
        return
    }

    // bmp header is 54 bytes
    pixelData := bmpData[54:]
    
    output := make([]byte, 800*480)
    
    // Convert RGBA to binary (using just the blue channel)
    for i := 0; i < 800*480; i++ {
        if pixelData[i*4] > 127 {
            output[i] = 0xFF  // white
        } else {
            output[i] = 0x00  // black
        }
    }

    w.Header().Set("Content-Type", "application/octet-stream")
    w.Write(output)
}

func main() {
  static, media, port = get_vars()
  http.Handle("/", http.FileServer(http.Dir(static)))

  http.HandleFunc("/get-bin", get_bin_handler)
  http.HandleFunc("/save", save_handler)
  http.HandleFunc("/get", get_handler)

  log.Println("Serving on port " + port)
  log.Println("media: " + media)
  err := http.ListenAndServe(":" + port, nil)
  if err != nil {
    panic(err)
  }
}
