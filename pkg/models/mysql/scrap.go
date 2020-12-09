package mysql

import (
	"database/sql"

	"github.com/gopheramit/web-scrapping/pkg/models"
)

type ScrapModel struct {
	DB *sql.DB
}

func (m *ScrapModel) Insert(email string) (int, error) {
	return 0, nil
}

// This will return a specific snippet based on its id.
func (m *ScrapModel) Get(id int) (*models.Scrap, error) {
	return nil, nil
}

// This will return the 10 most recently created snippets.
func (m *ScrapModel) Latest() ([]*models.Scrap, error) {
	return nil, nil
}
