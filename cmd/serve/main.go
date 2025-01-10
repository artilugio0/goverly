package main

import (
	"embed"
	"fmt"
	"net/http"
)

//go:embed wasm/*
var content embed.FS

func main() {
	http.Handle("/", http.FileServer(http.FS(content)))

	fmt.Println("Serving...")
	http.ListenAndServe(":8080", nil)
}
