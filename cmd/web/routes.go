package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/about", app.about)
	mux.HandleFunc("/getstarted", app.getStarted)
	fileServer := http.FileServer(http.Dir("./ui/static/main.css"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	return mux
}
