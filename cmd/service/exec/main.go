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
)

func main() {

	dsn := flag.String("dsn", "scrapit:pass@/webscrap?parseTime=true", "MySQL data source name")
	flag.Parse()

	// infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db1, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db1.Close()
	ql := service.NewQueueListener(db1)
	go ql.ListenForNewSource()
	//ql.ScrapRequest = DB: db
	var a string
	fmt.Scanln(&a)
}

func openDB(dsn string) (*sql.DB, error) {
	db1, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db1.Ping(); err != nil {
		return nil, err
	}
	return db1, nil
}
