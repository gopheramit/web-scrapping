package main

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/gopheramit/web-scrapping/pkg/forms"
	"github.com/gopheramit/web-scrapping/pkg/models"
	"github.com/markbates/goth/gothic"
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

func (app *application) pricing(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "pricing.page.tmpl", &templateData{
		//Scrap: s,
	})
	//w.Write([]byte("About pricing!"))

}

/*
func (app *application) createScarp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	email := "abcd@gmail.com"
	guid := "asdfghkl"
	expires := "8"

	id, err := app.scraps.Insert(email, "amit", guid, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/scrap?id=%d", id), http.StatusSeeOther)
}
*/
func (app *application) showScrap(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	//	var=r.URL.Query().Get("id")
	//	fmt.Println(var)
	//fmt.Println(id)
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	authenticatedUserID := app.session.Get(r, "authenticatedUserID")
	//fmt.Println("authenticatedUserID", authenticatedUserID)
	if authenticatedUserID == id {
		s, err := app.scraps.Get(id)
		//fmt.Println(s.Count)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				app.notFound(w)
			} else {
				app.serverError(w, err)
			}

			return
		}
		app.render(w, r, "show.page.tmpl", &templateData{
			Scrap: s,
		})

	} else {
		app.notFound(w)
		return
	}

	// Write the snippet data as a plain-text HTTP response body.
	//fmt.Fprintf(w, "%v", s)
}

func (app *application) authbegin(w http.ResponseWriter, r *http.Request) {
	fmt.Println("begignig authorisation!")
	gothic.BeginAuthHandler(w, r)
}
func (app *application) auth(w http.ResponseWriter, r *http.Request) {
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		fmt.Fprintln(w, r)
		return
	}
	t, _ := template.ParseFiles("ui/html/success.html")
	t.Execute(w, user)
}

///////////////////////////////////////////////////////////////////////////////////////////

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login1.page.tmpl", nil)
}

func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "signup1.page.tmpl", &templateData{
		// Pass a new empty forms.Form object to the template.
		Form: forms.New(nil),
	})
}

func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form := forms.New(r.PostForm)
	fmt.Println(form.Get("email"))
	form.Required("email", "password")
	form.MaxLength("email", 255)
	form.MatchesPattern("email", forms.EmailRX)
	form.MinLength("password", 2)
	if !form.Valid() {
		app.render(w, r, "signup1.page.tmpl", &templateData{
			// Pass a new empty forms.Form object to the template.
			Form: forms.New(nil),
		})
		return
	}
	//fmt.Println("amit")
	key := app.genUlid()
	keystr := key.String()
	count := 1000
	id, err := app.scraps.Insert(form.Get("email"), form.Get("password"), keystr, count, "30")
	//rr = app.users.Insert(form.Get("email"), form.Get("password"))
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.Errors.Add("email", "Address is already in use")
			app.render(w, r, "signup1.page.tmpl", nil)
		} else {

			app.serverError(w, err)
		}
		return
	}
	fmt.Println(id)
	// Otherwise send a placeholder response (for now!).
	//app.session.Put(r, "flash", "Your signup was successful. Please log in.")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)

}

func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}
func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// Check whether the credentials are valid. If they're not, add a generic error
	// message to the form failures map and re-display the login page.
	form := forms.New(r.PostForm)
	id, err := app.scraps.Authenticate(form.Get("email"), form.Get("password"))
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.Errors.Add("generic", "Email or Password is incorrect")
			app.render(w, r, "login.page.tmpl", &templateData{Form: form})
		} else {
			app.serverError(w, err)
		}
		return
	}
	// Add the ID of the current user to the session, so that they are now 'logged
	// in'.
	fmt.Println(id)

	app.session.Put(r, "authenticatedUserID", id)

	//fmt.Println(r)
	// Redirect the user to the create snippet page.
	//http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
	http.Redirect(w, r, fmt.Sprintf("/scrap/%d", id), http.StatusSeeOther)
}
func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	app.session.Remove(r, "authenticatedUserID")
	app.session.Put(r, "flash", "You've been logged out successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) linkScrape(w http.ResponseWriter, r *http.Request) {
	key := (r.URL.Query().Get(":key"))
	s, err := app.scraps.GetKey(key)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}

		return
	}
	fmt.Println(s.Count)
	if s.Count > 0 {
		res, err := http.Get("http://jonathanmh.com")
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()

		n, err := io.Copy(os.Stdout, res.Body)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Number of bytes copied to STDOUT:", n)

		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		doc.Find("body a").Each(func(index int, item *goquery.Selection) {
			linkTag := item
			link, _ := linkTag.Attr("href")
			linkText := linkTag.Text()
			fmt.Printf("Link #%d: '%s' - '%s'\n", index, linkText, link)
		})

		cnt := s.Count - 1
		fmt.Println(cnt)
		_, err = app.scraps.Decrement(s.ID, cnt)
		if err != nil {
			fmt.Println("error here")
		
		}

	} else {
		app.notFound(w)
		return
	}

	// doc, err := goquery.NewDocument("http://jonathanmh.com")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// print(doc)

}
