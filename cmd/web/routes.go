package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

//might need to change third party routing handler

func (app *application) routes() http.Handler {

	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeader)
	//mux := http.NewServeMux()
	mux := pat.New()
	mux.Get("/", http.HandlerFunc(app.home))
	mux.Get("/about", http.HandlerFunc(app.about))
	mux.Get("/documentation", http.HandlerFunc(app.documentation))
	mux.Get("/pricing", http.HandlerFunc(app.pricing))
	mux.Get("/login", http.HandlerFunc(app.login))
	mux.Get("/signup", http.HandlerFunc(app.signup))
	mux.Post("/showkeys", http.HandlerFunc(app.showkey))

	mux.Get("/scrap", http.HandlerFunc(app.showScrap))

	fileServer := http.FileServer(http.Dir("./assets/"))
	mux.Get("/assets/", http.StripPrefix("/assets/", fileServer))

	return standardMiddleware.Then(mux)
}
