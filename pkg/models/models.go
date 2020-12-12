package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: no matching record found")

type Scrap struct {
	ID      int
	Email   string
	Created time.Time
	Expires time.Time
}
