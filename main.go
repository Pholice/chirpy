package main

import (
	"net/http"
)

func main() {
	serveMux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./"))
	server := http.Server{
		Addr:    ":8080",
		Handler: serveMux,
	}
	serveMux.Handle("/", fileServer)
	server.ListenAndServe()
}
