package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/about", about)
	mux.HandleFunc("/getstarted", getStarted)

	log.Println("starting server on 8000")
	err := http.ListenAndServe(":8000", mux)
	log.Fatal(err)
}
