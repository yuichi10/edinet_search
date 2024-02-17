package db

import (
	"fmt"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

const dbName = "company.db"

type Documents struct {
	DocID               string `gorm:"primaryKey"`
	SecCode             string
	FilerName           string `gorm:"index"`
	DocDescription      string
	SubmitDatetime      string
	AvgAge              string
	AvgYearOfService    string
	AvgAnnualSalary     string `gorm:"index"`
	EmployeeInformation string
}

func DeleteDB() error {
	err := os.Remove(dbName)
	if err != nil {
		fmt.Println("DBの削除に失敗しました。 ", err)
	}
	return nil
}

func OpenDB() error {
	var err error
	db, err = gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("DBへの接続失敗 %s", err)
	}
	return nil
}

func CreateDocumentTable() error {
	err := db.AutoMigrate(&Documents{})
	if err != nil {
		return fmt.Errorf("テーブルの初期化に失敗 %s", err)
	}
	return nil
}

func InsertDocument(d Documents) error {
	result := db.Create(d)
	if result.Error != nil {
		return fmt.Errorf("ドキュメントデータの挿入に失敗 %s", result.Error)
	}
	return nil
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
