package main

import (
	"embed"
	"flag"
	"fmt"

	"github.com/artilugio0/goverly"
)

//go:embed wasm/*
var content embed.FS

func main() {
	configFile := flag.String("f", "config.json", "configuration file")
	flag.Parse()
	args := flag.Args()

	fmt.Println(*configFile)

	subcommand := ""
	if len(args) > 0 {
		subcommand = args[0]
	}

	switch subcommand {
	case "overlay":
		goverly.ServeOverlay(*configFile, content)
	case "config":
		fmt.Println("modifying config")
	}
}
