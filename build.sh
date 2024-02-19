#!/bin/bash
GOOS=js GOARCH=wasm go build -o dist/editor.wasm game-engine/editor.go game-engine/utils.go
goexec 'http.ListenAndServe(":8080", http.FileServer(http.Dir(".")))'
