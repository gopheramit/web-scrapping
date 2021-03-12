package main

import (
	//"distributed/coordinator"
	"fmt"

	"github.com/gopheramit/web-scrapping/cmd/service"
)

func main() {
	ql := service.NewQueueListener()
	go ql.ListenForNewSource()

	var a string
	fmt.Scanln(&a)
}
