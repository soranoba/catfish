package main

import (
	"log"
	"net/http"
)

func main() {
	handler := &HTTPHandler{}
	srv := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}
	log.Fatal(srv.ListenAndServe())
}
