package models

import (
	"errors"
)

var ErrNoRecord = errors.New("models: no matching record found")

type Scrap struct {
	ID    int
	Emial string
}
