package mysql

import (
	"database/sql"
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
