package main

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
)

//go:embed public/*
var publicFolderFS embed.FS

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Websocket Endpoint")
}

func main() {
    fmt.Println("yeah")

    publicContentFS, _ := fs.Sub(publicFolderFS, "public")

    mux := http.NewServeMux()
    mux.Handle("/", http.FileServer(http.FS(publicContentFS)))
    mux.HandleFunc("/ws", wsEndpoint)

    http.ListenAndServe(":3000", mux)
}
