package main

import (
	"net/http"

	"github.com/bmizerany/pat"

	//"github.com/go-chi/chi"
	"github.com/gorilla/sessions"
	"github.com/justinas/alice"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

//might need to change third party routing handler(using chi now)

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
		"379756554270-olm9ma6g4dru3lil2cse84eaeimpj0u2.apps.googleusercontent.com",
		"vlDxjmHJX80vOuHa5THxfCsR",
		"http://localhost:4000/auth/callback?provider=google", "email", "profile"))

	//mux := chi.NewRouter()
	mux := pat.New()
	mux.Get("/", app.session.Enable(http.HandlerFunc(app.home)))
	mux.Get("/about", http.HandlerFunc(app.about))
	mux.Get("/documentation", http.HandlerFunc(app.documentation))
	mux.Get("/pricing", http.HandlerFunc(app.pricing))
	mux.Get("/login", http.HandlerFunc(app.login))
	mux.Get("/auth/callback", http.HandlerFunc(app.auth))
	mux.Get("/auth", http.HandlerFunc(gothic.BeginAuthHandler))
	mux.Get("/signup", http.HandlerFunc(app.signupUserForm))
	//mux.Post("/signup", http.HandlerFunc(app.signup))
	mux.Get("/scrap/:id", http.HandlerFunc(app.showScrap))

	mux.Get("/user/signup", http.HandlerFunc(app.signupUserForm))
	mux.Post("/user/signup", http.HandlerFunc(app.signupUser))
	mux.Get("/user/login", http.HandlerFunc(app.loginUserForm))
	mux.Post("/user/login", http.HandlerFunc(app.loginUser))
	mux.Post("/user/logout", http.HandlerFunc(app.logoutUser))

	mux.Get("/snippet/create", http.HandlerFunc(app.createSnippetForm))

	mux.Post("/snippet/create", http.HandlerFunc(app.createSnippet))

	//filesDir := http.Dir("./assets/")
	//FileServer(mux, "/assets", filesDir)
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))
	return standardMiddleware.Then(mux)
}

/*
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
*/
