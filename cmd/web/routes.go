package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	//"github.com/gorilla/pat"
	"github.com/gorilla/sessions"
	"github.com/justinas/alice"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

//might need to change third party routing handler

func (app *application) routes() http.Handler {

	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeader)
	//mux := http.NewServeMux()
	key := "DoN4QZCXaa3TJfr4BJZMQZNo" // Replace with your SESSION_SECRET or similar
	maxAge := 86400 * 30              // 30 days
	isProd := false                   // Set to true when serving over https

	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(maxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true // HttpOnly should always be enabled
	store.Options.Secure = isProd

	gothic.Store = store
	goth.UseProviders(
		google.New("263741611747-2bgmmh2vnbjvt02c3m8s30ujbb76obgf.apps.googleusercontent.com", "DoN4QZCXaa3TJfr4BJZMQZNo", "http://localhost:4000/auth/google/callback", "email", "profile"),
	)
	mux := pat.New()
	mux.Get("/auth/{provider}/callback", http.HandlerFunc(app.Auth))
	mux.Get("/auth/{provider}", http.HandlerFunc(app.authbegin))
	mux.Get("/", http.HandlerFunc(app.home))
	mux.Get("/about", http.HandlerFunc(app.about))
	mux.Get("/documentation", http.HandlerFunc(app.documentation))
	mux.Get("/pricing", http.HandlerFunc(app.pricing))
	mux.Get("/login", http.HandlerFunc(app.login))
	mux.Get("/signup", http.HandlerFunc(app.signupForm))
	mux.Post("/signup", http.HandlerFunc(app.signup))

	mux.Get("/scrap/:id", http.HandlerFunc(app.showScrap))

	fileServer := http.FileServer(http.Dir("./assets/"))
	mux.Get("/assets/", (http.StripPrefix("/assets/", fileServer)))

	return standardMiddleware.Then(mux)
}
