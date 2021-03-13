package main

import (
	//"distributed/coordinator"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	//"github.com/gopheramit/web-scrapping/cmd/service"
	"github.com/gopheramit/web-scrapping/cmd/service"
	"github.com/gopheramit/web-scrapping/cmd/service/models"
)

type application1 struct {
	ScrapRequest *models.ScrapRequestModel
}

func main() {
	ql := service.NewQueueListener()
	go ql.ListenForNewSource()
	dsn := flag.String("dsn", "scrapit:pass@/webscrap?parseTime=true", "MySQL data source name")
	flag.Parse()

	// infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()
	_ = &application1{
		ScrapRequest: &models.ScrapRequestModel{DB: db},
	}

	var a string
	fmt.Scanln(&a)
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
