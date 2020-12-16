package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"

	//"github.com/gorilla/pat"
	"github.com/go-chi/chi"
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
	key := "vlDxjmHJX80vOuHa5THxfCsR" // Replace with your SESSION_SECRET or similar
	maxAge := 86400 * 30              // 30 days
	isProd := false                   // Set to true when serving over https

	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(maxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true // HttpOnly should always be enabled
	store.Options.Secure = isProd

	gothic.Store = store
	goth.UseProviders(google.New(
		os.Getenv("379756554270-olm9ma6g4dru3lil2cse84eaeimpj0u2.apps.googleusercontent.com"),
		os.Getenv("vlDxjmHJX80vOuHa5THxfCsR"),
		"http://localhost:4000/auth/callback?provider=google", "email", "profile"))

	//mux := pat.New()
	mux := chi.NewRouter()
	mux.Get("/", http.HandlerFunc(app.home))
	mux.Get("/about", http.HandlerFunc(app.about))
	mux.Get("/documentation", http.HandlerFunc(app.documentation))
	mux.Get("/pricing", http.HandlerFunc(app.pricing))
	mux.Get("/login", http.HandlerFunc(app.login))
	//mux.Get("auth/{provider}/callback", http.HandlerFunc(app.auth))
	mux.Get("/auth/callback", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := gothic.CompleteUserAuth(w, r)
		if err != nil {
			fmt.Fprintln(w, r)
			return
		}
		t, _ := template.ParseFiles("ui/html/success.html")
		t.Execute(w, user)
	}))

	mux.Get("/auth", gothic.BeginAuthHandler)
	/*mux.Get("/auth/google", http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		gothic.BeginAuthHandler(res, req)
	}))*/
	mux.Get("/signup", http.HandlerFunc(app.signupForm))
	mux.Post("/signup", http.HandlerFunc(app.signup))

	mux.Get("/scrap/:id", http.HandlerFunc(app.showScrap))
	filesDir := http.Dir("./assets/")
	FileServer(mux, "/assets", filesDir)

	//fileServer := http.FileServer(http.Dir("./assets/"))
	//http.StripPrefix("/assets/", fileServer)

	//mux.Get("/assets/", http.Handle(myHandler))

	return standardMiddleware.Then(mux)
}

func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
