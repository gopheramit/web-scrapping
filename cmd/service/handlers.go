package service

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/go-sql-driver/mysql"
	"github.com/oklog/ulid"
)

func genUlid() ulid.ULID {
	t := time.Now().UTC()
	entropy := rand.New(rand.NewSource(t.UnixNano()))
	id := ulid.MustNew(ulid.Timestamp(t), entropy)
	return id
}

//add swagger for following handler.

func (ql *QueueListener) linkscrape(url string) { // key string) {
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

	key1 := genUlid()
	keystr1 := key1.String()
	fmt.Println(keystr1)
	resullt, err := doc.Html()

	//var final []byte
	//	final = []byte(resullt)
	err = ql.ScrapRequest.Insert(keystr1, keystr1, []byte(resullt))
	if err != nil {
		fmt.Println("error linkscrape")

	} else {
		fmt.Println("everthing ok")
	}
}

// buf := new(bytes.Buffer)
// enc := gob.NewEncoder(buf)

// enc.Encode(doc.Html)
/////////////////////////////////////////////////////////////////////////////////////////
/*
		stmt := `INSERT INTO bloob(BLOBdata,created,count)VALUES(?,?,UTC_TIMESTAMP(),?)`

		_, err = db.Exec(stmt, doc, 1000)

		if err != nil {
			fmt.Println(err)
		}

//////////////////////////////////////////////////////////////////////////////////////////
//return doc.Html()
//resullt, err := doc.Html()
//w.Write([]byte(resullt))

/*
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
*/
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
		//**
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
