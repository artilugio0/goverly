#!/usr/bin/env bash

if [[ ! -d ./build ]] then
	mkdir ./build
fi

# compile WASM binary
GOOS=js GOARCH=wasm go build -o ./build/main.wasm ./cmd/overlay/
mv ./build/main.wasm ./cmd/goverly/wasm/main.wasm

# compile goverly binary
go build -o ./build/goverly ./cmd/goverly
mv ./cmd/goverly/wasm/main.wasm ./build/main.wasm
