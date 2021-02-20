package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/gopheramit/web-scrapping/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

type ScrapModel struct {
	DB *sql.DB
}

func (m *ScrapModel) Insert(socID, email, password, guid string, count int, expires string) (int, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)

	if err != nil {
		return 0, err
	}
	stmt := `INSERT INTO scraps (soc_id,email, hashed_password,guid,count,created, expires)VALUES(?,?,?,?, ?,UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`
	result, err := m.DB.Exec(stmt, socID, email, string(hashedPassword), guid, count, expires)
	if err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return 0, models.ErrDuplicateEmail
			}
		}
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

	stmt := `SELECT id,soc_id, email,guid,count,created, expires FROM scraps WHERE expires > UTC_TIMESTAMP() AND id = ?`
	row := m.DB.QueryRow(stmt, id)
	s := &models.Scrap{}
	err := row.Scan(&s.ID, &s.Soc_id, &s.Email, &s.Guid, &s.Count, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil

}

func (m *ScrapModel) GetKey(id string) (*models.Scrap, error) {

	stmt := `SELECT id, email,guid,count,created, expires FROM scraps WHERE expires > UTC_TIMESTAMP() AND guid= ?`
	row := m.DB.QueryRow(stmt, id)
	s := &models.Scrap{}

	err := row.Scan(&s.ID, &s.Email, &s.Guid, &s.Count, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil

}

func (m *ScrapModel) Decrement(id, count int) (int, error) {
	fmt.Println(id, count)
	//stmt := `SELECT id, email,guid,created, expires FROM scraps WHERE expires > UTC_TIMESTAMP() AND guid= ?`
	stmt := `update scraps set count=? where id=?`
	_, err := m.DB.Exec(stmt, count, id)
	if err != nil {
		return 1, err
	}
	return 0, nil
}

// This will return the 10 most recently created snippets.
func (m *ScrapModel) Latest() ([]*models.Scrap, error) {
	return nil, nil
}
func (m *ScrapModel) GetID(socID string) (*models.Scrap, error) {
	//fmt.Println("in GEtID")
	stmt := `SELECT id, soc_id,email,guid,count,created, expires FROM scraps WHERE expires > UTC_TIMESTAMP() AND soc_id=?`
	row := m.DB.QueryRow(stmt, socID)
	s := &models.Scrap{}

	err := row.Scan(&s.ID, &s.Soc_id, &s.Email, &s.Guid, &s.Count, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}

func (m *ScrapModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte
	stmt := "SELECT id, hashed_password FROM scraps WHERE email = ? AND active = TRUE"
	row := m.DB.QueryRow(stmt, email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	// Check whether the hashed password and plain-text password provided match.
	// If they don't, we return the ErrInvalidCredentials error.
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	// Otherwise, the password is correct. Return the user ID.
	return id, nil
}
