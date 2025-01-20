package main

import (
	"embed"
	"fmt"
	"os"

	"github.com/artilugio0/goverly"
)

//go:embed wasm/*
var content embed.FS

func main() {
	subcommand := ""
	if len(os.Args) >= 2 {
		subcommand = os.Args[1]
	}

	switch subcommand {
	case "overlay":
		goverly.ServeOverlay("config.json", content)
	case "config":
		fmt.Println("modifying config")
	}
}
