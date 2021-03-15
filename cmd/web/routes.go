package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/gorilla/sessions"

	//"github.com/go-chi/chi"
	// "github.com/gorilla/pat"
	"github.com/justinas/alice"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

//might need to change third party routing handler(using chi now)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeader)
	dynamicMiddleware := alice.New(app.session.Enable) //, noSurf)
	//mux := http.NewServeMux()
	key1 := "DoN4QZCXaa3TJfr4BJZMQZNo" // Replace with your SESSION_SECRET or similar
	maxAge := 86400 * 30               // 30 days
	isProd := false                    // Set to true when serving over https
	store := sessions.NewCookieStore([]byte(key1))
	store.MaxAge(maxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true // HttpOnly should always be enabled
	store.Options.Secure = isProd
	gothic.Store = store
	// goth.UseProviders(google.New(
	// 	"379756554270-olm9ma6g4dru3lil2cse84eaeimpj0u2.apps.googleusercontent.com",
	// 	"vlDxjmHJX80vOuHa5THxfCsR",
	// 	"http://localhost:4000/auth/callback?provider=google", "email", "profile"))
	goth.UseProviders(
		google.New("263741611747-2bgmmh2vnbjvt02c3m8s30ujbb76obgf.apps.googleusercontent.com", "DoN4QZCXaa3TJfr4BJZMQZNo", "http://localhost:4000/auth/google/callback", "email", "profile"),
	)

	//mux := chi.NewRouter()
	mux := pat.New()
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	mux.Get("/about", dynamicMiddleware.ThenFunc(app.about))
	mux.Get("/documentation", dynamicMiddleware.ThenFunc(app.documentation))
	mux.Get("/pricing", dynamicMiddleware.ThenFunc(app.pricing))
	//mux.Get("/login", dynamicMiddleware.ThenFunc(app.login))
	mux.Get("/auth/:provider/callback", dynamicMiddleware.ThenFunc(app.auth))
	mux.Get("/auth/:provider", dynamicMiddleware.ThenFunc(gothic.BeginAuthHandler))
	mux.Get("/scrap/:id", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.showScrap))

	mux.Get("/user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm))
	mux.Post("/user/signup", dynamicMiddleware.ThenFunc(app.signupUser))
	mux.Get("/user/verify", dynamicMiddleware.ThenFunc(app.VerifyUserForm))
	mux.Post("/user/verify", dynamicMiddleware.ThenFunc(app.VerifyUser))
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.loginUser))
	mux.Get("/user/resend", dynamicMiddleware.ThenFunc(app.resendOtp))
	mux.Post("/user/logout", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.logoutUser))

	//mux.Get("/request", http.HandlerFunc(app.linkScrape))
	//mux.Get("/requestheaders", http.HandlerFunc(app.linkScrapeheaders))
	//mux.Get("/request/render", http.HandlerFunc(app.JsRendering))
	mux.Get("/request", http.HandlerFunc(app.Decision))
	mux.Get("/requestheaders", http.HandlerFunc(app.Decision))
	mux.Get("/request/render", http.HandlerFunc(app.Decision))
	//mux.Get("/echo", http.HandlerFunc(app.echo))
	http.HandleFunc("/echo", app.echo)
	//filesDir := http.Dir("./assets/")
	//FileServer(mux, "/assets", filesDir)
	fileServer := http.FileServer(http.Dir("./assets"))
	mux.Get("/assets/", http.StripPrefix("/assets", fileServer))
	return standardMiddleware.Then(mux)
}
