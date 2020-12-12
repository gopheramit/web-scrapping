package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/gopheramit/web-scrapping/pkg/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}
	//data := &templateData{Scrap: s}
	/*files := []string{
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

	err = ts.Execute(w, data)
	if err != nil {
		app.serverError(w, err)
	}*/
	s, err := app.scraps.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.render(w, r, "home.page.tmpl", &templateData{
		Scraps: s,
	})
	//w.Write([]byte("hello from scrapper!"))

}

func (app *application) about(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("About webscarpping!"))
}

func (app *application) documentation(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Get started with web scrapping!"))
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.tmpl", nil)
}

func (app *application) pricing(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "pricing.page.tmpl", nil)
	//w.Write([]byte("About pricing!"))

}

func (app *application) createScarp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	email := "abcd@gmail.com"
	guid := "asdfghkl"
	expires := "8"

	id, err := app.scraps.Insert(email, guid, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/scrap?id=%d", id), http.StatusSeeOther)
}

func (app *application) showScrap(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	s, err := app.scraps.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}

		return
	}
	data := &templateData{Scrap: s}
	files := []string{
		"./ui/html/show.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}
	err = ts.Execute(w, data)
	if err != nil {
		app.serverError(w, err)
	}
	// Write the snippet data as a plain-text HTTP response body.
	fmt.Fprintf(w, "%v", s)
}

func (app *application) signupForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "signup.page.tmpl", nil)
}

func (app *application) signup(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	errors := make(map[string]string)
	email := r.PostForm.Get("Email")
	fmt.Println(email)
	if strings.TrimSpace(email) == "" {
		errors["email"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(email) > 100 {
		errors["email"] = "This field is too long (maximum is 100 characters)"
	}
	if len(errors) > 0 {
		app.render(w, r, "signup.page.tmpl", &templateData{
			FormErrors: errors,
			FormData:   r.PostForm,
		})
		return
	}
	key := app.genUlid()
	keystr := key.String()
	id, err := app.scraps.Insert(email, keystr, "30")
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/scrap/%d", id), http.StatusSeeOther)
	//app.render(w, r, "keys.page.tmpl", nil)
}
