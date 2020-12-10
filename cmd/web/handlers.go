package main

import (
	"net/http"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}
	/*
		files := []string{
			//"./ui/html/index.html",
			"./ui/html/home.page.tmpl",
			"./ui/html/base.layout.tmpl",
			"./ui/html/footer.partial.tmpl",
		}

		ts, err := template.ParseFiles(files...)
		if err != nil {
			app.serverError(w, err)
			return
		}

		err = ts.Execute(w, nil)
		if err != nil {
			app.serverError(w, err)
		}*/
	app.render(w, r, "home.page.tmpl")
	//w.Write([]byte("hello from scrapper!"))

}

func (app *application) about(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("About webscarpping!"))
}

func (app *application) documentation(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Get started with web scrapping!"))
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.tmpl")
}

func (app *application) pricing(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "pricing.page.tmpl")
	//w.Write([]byte("About pricing!"))

}

func (app *application) signup(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "signup.page.tmpl")
}

func (app *application) showkey(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "keys.page.tmpl")
}
