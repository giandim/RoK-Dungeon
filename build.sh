#!/bin/bash

# Function to kill the server
kill_server() {
	echo "Killing server..."
	# Replace the following line with the command to kill your server
	kill -TERM $SERVER_PID
}

# Set up trap to call kill_server function on EXIT and TERM signals
trap kill_server EXIT TERM

# Build Go WebAssembly module
GOOS=js GOARCH=wasm go build -o dist/editor.wasm game-engine/*.go

# Build Go server executable
go build -o dist/server server/server.go

# Start your server in the background and store its PID
./dist/server &
SERVER_PID=$!

# Start the Go WebAssembly server using goexec
goexec 'http.ListenAndServe(":8080", http.FileServer(http.Dir(".")))'

# Wait for the server process to finish (if it finishes before the build completes)
wait $SERVER_PID
