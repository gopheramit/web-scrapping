package mysql

/*
import (
	"database/sql"
	"errors" // New import
	"strings"

	"github.com/go-sql-driver/mysql" // New import
	"github.com/gopheramit/web-scrapping/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(email, password string) error {
	// Create a bcrypt hash of the plain-text password.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	stmt := `INSERT INTO users (email, hashed_password, created)VALUES(?, ?, UTC_TIMESTAMP())`

	_, err = m.DB.Exec(stmt, email, string(hashedPassword))
	if err != nil {

		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return models.ErrDuplicateEmail
			}
		}
		return err
	}
	return nil
}
*/
