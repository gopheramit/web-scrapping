package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gopheramit/web-scrapping/pkg/models/mysql"
	"github.com/gorilla/websocket"

	//"github.com/gorilla/sessions"
	"github.com/golangcollege/sessions"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	templateCache map[string]*template.Template
	scraps        *mysql.ScrapModel
	otps          *mysql.OtpModel
	session       *sessions.Session
	Key           *string
	ScrapRequest  *mysql.ScrapRequestModel
	//ws            *Websocket
}

//type contextKey string

//const contextKeyIsAuthenticated = contextKey("isAuthenticated")
var upgrader = websocket.Upgrader{}

func main() {

	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "scrapit:pass@/webscrap?parseTime=true", "MySQL data source name")
	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret key")
	usrKey := flag.String("key", " 01ETWM58TWCWJ3JZYWH2Q33B1N", "UserKey")
	flag.Parse()
	session := sessions.New([]byte(*secret))
	session.Lifetime = 24 * time.Hour
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	templateCache, err := newTemplateCache("./ui/html/")

	//ws := new WebSocket("ws://localhost:4000/")
	if err != nil {
		errorLog.Fatal(err)
	}

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		templateCache: templateCache,
		scraps:        &mysql.ScrapModel{DB: db},
		otps:          &mysql.OtpModel{DB: db},
		session:       session,
		Key:           usrKey,
		ScrapRequest:  &mysql.ScrapRequestModel{DB: db},
		//ws:            Websocket,
	}
	//usrKey := "01ETWM58TWCWJ3JZYWH2Q33B1N"

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}
	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
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
