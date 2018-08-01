package main

import (
	"log"
	"net/http"
	"os"
)

type myFileSystem struct {
	http.FileSystem
}

func main() {
	home := os.Getenv("HOME")
	customFileSystem := myFileSystem{http.Dir(home)}
	http.Handle("/", http.FileServer(customFileSystem))
	log.Fatal(http.ListenAndServe(":12346", nil))
}
