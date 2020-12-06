package main

import (
	"net/http"

	"github.com/justinas/alice"
)

//might need to change third party routing handler

func (app *application) routes() http.Handler {

	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeader)
	mux := http.NewServeMux()
	go mux.HandleFunc("/", app.home)
	mux.HandleFunc("/about", app.about)
	mux.HandleFunc("/documentation", app.documentation)
	mux.HandleFunc("/pricing", app.pricing)

	fileServer := http.FileServer(http.Dir("./assets/"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fileServer))

	return standardMiddleware.Then(mux)
}
