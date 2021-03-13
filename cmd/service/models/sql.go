package mysql

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/gopheramit/web-scrapping/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

type ScrapRequestModel struct {
	DB *sql.DB
}

func (m *ScrapRequestModel) Insert() (int, error) {


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
