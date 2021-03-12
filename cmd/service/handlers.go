package service

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

//add swagger for following handler.
func linkscrape(url string) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(doc.Html())
	//return doc.Html()
	//resullt, err := doc.Html()
	//w.Write([]byte(resullt))
}

/*
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
*/
/*
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
