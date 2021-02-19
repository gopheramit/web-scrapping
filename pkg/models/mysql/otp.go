package mysql

import (
	"database/sql"
	"errors"

	"github.com/gopheramit/web-scrapping/pkg/models"
)

type OtpModel struct {
	DB *sql.DB
}

func (m *OtpModel) InsertOtp(id int, otp string) error {

	stmt := `INSERT INTO Otps (id,otp,created, expires)VALUES(?,?,UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`
	_, err := m.DB.Exec(stmt, id, otp, "1")
	if err != nil {
		return err
	}
	return nil
}

func (m *OtpModel) GetOtp(id int) (*models.Otps, error) {

	stmt := `SELECT id,otp,created, expires FROM scraps WHERE expires > UTC_TIMESTAMP() AND id = ?`
	row := m.DB.QueryRow(stmt, id)
	s := &models.Otps{}
	err := row.Scan(&s.ID, &s.Otp, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil

}
