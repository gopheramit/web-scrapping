package main

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"net/smtp"
	"strconv"

	"github.com/gopheramit/web-scrapping/pkg/forms"
	"github.com/gopheramit/web-scrapping/pkg/models"
	"github.com/markbates/goth/gothic"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}
	s, err := app.scraps.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.render(w, r, "home.page.tmpl", &templateData{
		Scraps: s,
	})
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

}

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
	//fmt.Println(user.Email)
	//fmt.Println(user.UserID)
	if err != nil {
		fmt.Fprintln(w, r)
		fmt.Println("error here")
		return
	}
	s, err := app.scraps.GetID(user.UserID)
	//fmt.Println(err)
	//fmt.Println(s)
	if s != nil {
		//	fmt.Println("user found", s.ID)
		app.session.Put(r, "authenticatedUserID", s.ID)
		http.Redirect(w, r, fmt.Sprintf("/scrap/%d", s.ID), http.StatusSeeOther)
		return
	} else {
		//fmt.Println("In else of linkscrape")
		key := app.genUlid()
		keystr := key.String()
		count := 1000
		form := forms.New(r.PostForm)
		id, err := app.scraps.Insert(user.UserID, user.Email, user.UserID, keystr, count, "30")
		if err != nil {
			if errors.Is(err, models.ErrDuplicateEmail) {
				form.Errors.Add("email", "Address is already in use")
				app.render(w, r, "login.page.tmpl", &templateData{Form: form})
			} else {

				app.serverError(w, err)
			}
			return
		}
		fmt.Println(id)
		app.session.Put(r, "authenticatedUserID", id)
		http.Redirect(w, r, fmt.Sprintf("/scrap/%d", id), http.StatusSeeOther)
		//t, _ := template.ParseFiles("ui/html/success.html")
		//fort.Execute(w, user)
	}
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
	socID := "1"
	id, err := app.scraps.Insert(socID, form.Get("email"), form.Get("password"), keystr, count, "30")
	//rr = app.users.Insert(form.Get("email"), form.Get("password"))
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.Errors.Add("email", "Address is already in use")
			app.render(w, r, "signup1.page.tmpl", &templateData{Form: form})
		} else {

			app.serverError(w, err)
		}
		return
	}

	otp, err := GenerateOTP(6)
	//userID = id
	from := "flutterproject13@gmail.com"
	password := "iskcon123"
	// Receiver email address.
	to := []string{
		"amittest53@gmail.com",
	}
	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)
	t, _ := template.ParseFiles("ui/html/verificationcode.html")
	var body bytes.Buffer
	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: This is a test subject \n%s\n\n", mimeHeaders)))
	t.Execute(&body, struct {
		Name    string
		Message string
	}{
		Name:    "Achal Agrawal",
		Message: "Your OTP is : " + otp,
	})
	// Sending email.
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Email Sent!")
	fmt.Println(otp)
	fmt.Println(id)
	err = app.otps.InsertOtp(id, otp)
	if err != nil {
		fmt.Println("error in otp insert ")

	}
	// Otherwise send a placeholder response (for now!).
	//app.session.Put(r, "flash", "Your signup was successful. Please log in.")
	app.session.Put(r, "UserID", id)
	http.Redirect(w, r, "/user/verify", http.StatusSeeOther)

}

func (app *application) resendOtp(w http.ResponseWriter, r *http.Request) {

	otp, err := GenerateOTP(6)
	//userID = id
	from := "flutterproject13@gmail.com"
	password := "iskcon123"
	// Receiver email address.
	to := []string{
		"amittest53@gmail.com",
	}
	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)
	t, _ := template.ParseFiles("ui/html/verificationcode.html")
	var body bytes.Buffer
	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: This is a test subject \n%s\n\n", mimeHeaders)))
	t.Execute(&body, struct {
		Name    string
		Message string
	}{
		Name:    "Achal Agrawal",
		Message: "Your OTP is : " + otp,
	})
	// Sending email.
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Email Sent!")
	fmt.Println(otp)
	userID, ok := app.session.Get(r, "UserID").(int)
	if !ok {
		app.serverError(w, errors.New("type assertion to string failed"))
	}
	fmt.Println(userID)
	_, err = app.otps.UpdateOtp(userID, otp)
	if err != nil {
		fmt.Println("error in otp insert ")

	}
	http.Redirect(w, r, "/user/verify", http.StatusSeeOther)
}

func (app *application) VerifyUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "verification.page.tmpl", &templateData{
		Form: forms.New(nil),
	})

}
func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) VerifyUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form := forms.New(r.PostForm)
	otp := form.Get("otp")
	fmt.Println("otp:", otp)
	userID, ok := app.session.Get(r, "UserID").(int)
	if !ok {
		app.serverError(w, errors.New("type assertion to string failed"))
	}
	s, err := app.otps.GetData(userID)
	if err != nil {
		fmt.Println(err)
		fmt.Println("error in verify user getotp")
	}

	if s.Otp == otp {
		_, err := app.otps.UppdateVerifyStatus(userID)
		if err != nil {
			//app.render(w, r, "verification.page.tmpl", nil)
			fmt.Println("Error in updating status")
		}
		s, err = app.otps.GetData(userID)
		if err != nil {
			fmt.Println(err)
			fmt.Println("error in verify user in getdata ")
		}
		if s.Verified == true {
			app.session.Put(r, "authenticatedUserID", userID)
			http.Redirect(w, r, fmt.Sprintf("/scrap/%d", userID), http.StatusSeeOther)
		} else {
			app.render(w, r, "verification.page.tmpl", &templateData{Form: forms.New(nil)})
			fmt.Println("Not verified")
		}
	} else {
		app.render(w, r, "verification.page.tmpl", &templateData{Form: forms.New(nil)})
		fmt.Println("Wrong Otp")
	}
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
	app.session.Put(r, "UserID", id)
	s, err := app.otps.GetData(id)
	if err != nil {
		fmt.Println("error while verifying user login")
	}
	if s.Verified == true {
		fmt.Println(id)
		app.session.Put(r, "authenticatedUserID", id)
		//fmt.Println(r)
		// Redirect the user to the create snippet page.
		//http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
		http.Redirect(w, r, fmt.Sprintf("/scrap/%d", id), http.StatusSeeOther)
	} else {
		app.render(w, r, "verification.page.tmpl", &templateData{Form: form})
		fmt.Println("User not verified reddirecting to verification page")
	}

}
func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	app.session.Remove(r, "authenticatedUserID")
	app.session.Put(r, "flash", "You've been logged out successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
func (app *application) Decision(w http.ResponseWriter, r *http.Request) {
	key := (r.URL.Query().Get("api_key"))
	url := r.URL.Query().Get("url")
	html := true
	js1 := r.URL.Query().Get("js")
	header := r.URL.Query().Get("header")
	fmt.Println(key, url, js1, header, html)
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
		main1(url,js1)
	}

}

/*
//add swagger for following handler.
func (app *application) linkScrape(w http.ResponseWriter, r *http.Request) {
	key := (r.URL.Query().Get("api_key"))
	url := r.URL.Query().Get("url")
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
		//res, err := http.Get("http://jonathanmh.com")
		res, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()

		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Println(doc.Html())
		//return doc.Html()
		resullt, err := doc.Html()
		w.Write([]byte(resullt))
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

}

func (app *application) linkScrapeheaders(w http.ResponseWriter, r *http.Request) {
	key := (r.URL.Query().Get("api_key"))
	url := r.URL.Query().Get("url")
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
		//****
		//  Use this url to test : http://httpbin.org/anything
		//**
		myCookie := &http.Cookie{
			Name:  "cookieKey1",
			Value: "value1",
		}

		req, _ := http.NewRequest("GET", url, nil)

		req.AddCookie(myCookie)

		req.Header.Add("x-rapidapi-key", "4fa6109f53msh9537939930788bap193d31jsn321563aa93d4")
		req.Header.Add("x-rapidapi-host", "scrapingbee.p.rapidapi.com")

		fmt.Println(req.Cookies())
		fmt.Println(req.Header)
		// req.Header.Add("Accept", "application/json")

		res, _ := http.DefaultClient.Do(req)

		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()

		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			log.Fatal(err)
		}

		resullt, err := doc.Html()
		w.Write([]byte(resullt))
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

}

func (app *application) JsRendering(w http.ResponseWriter, r *http.Request) {
	geziyor.NewGeziyor(&geziyor.Options{
		StartRequestsFunc: func(g *geziyor.Geziyor) {
			g.GetRendered("https://jonathanmh.com", g.Opt.ParseFunc)
		},
		ParseFunc: func(g *geziyor.Geziyor, r *client.Response) {
			resullt := r.HTMLDoc
			result, err := resullt.Html()
			if err != nil {
				fmt.Println("errorin js rendering handler")
			}
			w.Write([]byte(result))
			//fmt.Println(string(r.Body))
		},
		//BrowserEndpoint: "ws://localhost:8080",
	}).Start()
}
*/
