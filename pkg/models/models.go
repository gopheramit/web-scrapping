package models

import (
	"errors"
	"time"

	"github.com/oklog/ulid"
)

var ErrNoRecord = errors.New("models: no matching record found")

type Scrap struct {
	ID      int
	Email   string
	Guid    ulid.ULID
	Created time.Time
	Expires time.Time
}
