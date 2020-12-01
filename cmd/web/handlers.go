package main

import (
	"html/template"
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "internal server error", 500)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "internal server error", 500)
	}
	//w.Write([]byte("hello from scrapper!"))

}

func about(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("About webscarpping!"))
}

func getStarted(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Get started with web scrapping!"))
}
