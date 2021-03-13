package models

import (
	"database/sql"
	"fmt"
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
