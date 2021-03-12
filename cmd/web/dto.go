package main

import (
	"encoding/gob"
	"time"
)

type SensorMessage struct {
	Name      string
	Url       string
	Js        string
	Header    bool
	Html      bool
	Timestamp time.Time
}

func init() {
	gob.Register(SensorMessage{})
}
