package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// どこかでエラー処理ちゃんといれる。

var db *gorm.DB

const dbName = "company.db"

type Documents struct {
	DocID string	`gorm:"primaryKey"`
	SecCode string
	FilerName string `gorm:"index"`
	DocDescription string
	SubmitDatetime string
	AvgAge string
	AvgYearOfService string
	AvgAnnualSalary string `gorm:"index"`
	EmployeeInformation string
}

func DeleteDB() {
	os.Remove(dbName)
}

func OpenDB() {
	var err error
	db, err = gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		log.Fatal("DBへの接続失敗", err)
	}
}

func CreateDocumentTable() {
	db.AutoMigrate(&Documents{})
}

func InsertDocument(d Documents) {
	db.Create(d)
}

func GetDocuments(filerNames []string) ([]Documents, error) {
	var docs []Documents

	query := db.Model(&Documents{})
	for _, filerName := range filerNames {
		query.Or("filer_name LIKE ?", fmt.Sprintf("%%%s%%", filerName))
	}

	result := query.Find(&docs)

	if result.Error != nil {
		return nil, fmt.Errorf("ドキュメントの一覧取得失敗 %s", result.Error)
	}
	return docs, nil
}
