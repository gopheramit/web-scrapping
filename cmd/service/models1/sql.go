package models1

/*
import (
	"database/sql"
	"errors"
	"fmt"
	//"github.com/gopheramit/web-scrapping/pkg/models"
	//"github.com/gopheramit/web-scrapping/pkg/models"
)

type ScrapRequestModel struct {
	DB *sql.DB
}

func (m *ScrapRequestModel) Insert(uuid, guid string, BLOBData []byte) error {
	fmt.Println("insert in scraprequest")
	//fmt.Println(BLOBData)
	stmt := `INSERT INTO ScrapRequest (uuid,guid,BLOBData) VALUES (?,?,?)`
	_, err := m.DB.Exec(stmt, uuid, guid, BLOBData)
	if err != nil {
		fmt.Println("Error in adding scraped data in database")

		return err
	}
	return nil
}

func (m *ScrapRequestModel) GetData(guid string) (*ScrapRequest, error) {

	stmt := `SELECT BLOBData FROM ScrapRequest WHERE  guid = ?`
	row := m.DB.QueryRow(stmt, guid)
	s := &ScrapRequest{}
	err := row.Scan(&s.Uuid, &s.Guid, &s.BLOBData)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil

}
*/
