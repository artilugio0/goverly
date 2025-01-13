package main

import (
	"embed"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

//go:embed wasm/*
var content embed.FS

func main() {
	startTime := time.Now()
	http.Handle("/", http.FileServer(http.FS(content)))
	http.HandleFunc("/last-update", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(strconv.FormatInt(startTime.Unix(), 10)))
	})

	fmt.Println("Serving...")
	http.ListenAndServe(":8080", nil)
}
