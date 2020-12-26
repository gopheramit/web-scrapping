package mysql

import (
	"database/sql"
	"errors"

	"github.com/gopheramit/web-scrapping/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

type ScrapModel struct {
	DB *sql.DB
}

func (m *ScrapModel) Insert(email, password, guid, expires string) (int, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)

	if err != nil {
		return 0, err
	}
	stmt := `INSERT INTO scraps (email, hashed_password,guid,created, expires)VALUES(?,?,?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`
	result, err := m.DB.Exec(stmt, email, string(hashedPassword), guid, expires)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

// This will return a specific snippet based on its id.
func (m *ScrapModel) Get(id int) (*models.Scrap, error) {

	stmt := `SELECT id, email,guid,created, expires FROM scraps WHERE expires > UTC_TIMESTAMP() AND id = ?`
	row := m.DB.QueryRow(stmt, id)
	s := &models.Scrap{}
	err := row.Scan(&s.ID, &s.Email, &s.Guid, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil

}

// This will return the 10 most recently created snippets.
func (m *ScrapModel) Latest() ([]*models.Scrap, error) {
	return nil, nil
}
