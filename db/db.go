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
	DocID               string `gorm:"primaryKey"` // ドキュメントID
	SecCode             string // コード
	FilerName           string `gorm:"index"` // 上げた人・会社
	DocDescription      string // 見ている資料の情報
	SubmitDatetime      string // 更新された時期
	AvgAge              string // 働いている人の平均年齢
	AvgYearOfService    string // 働いている人が大体何年間働いているか
	AvgAnnualSalary     string `gorm:"index"` // 平均年収
	EmployeeInformation string // 授業員情報
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

func GetDocuments(filerNames []string, salary string) ([]Documents, error) {
	var docs []Documents

	query := db.Model(&Documents{})
	setWhere := false
	if len(filerNames) > 0 {
		if !setWhere {
			query.Where("filer_name LIKE ?", fmt.Sprintf("%%%s%%", filerNames[0]))
		}
		for i, filerName := range filerNames {
			if i == 0 && !setWhere {
				setWhere = true
				continue
			}
			query.Or("filer_name LIKE ?", fmt.Sprintf("%%%s%%", filerName))
		}
	}
	if salary != "" {
		if !setWhere {
			query.Where("avg_annual_salary > ?", salary)
			setWhere = true
		} else {
			query.Or("avg_annual_salary > ?", salary)
		}
	}

	result := query.Find(&docs)

	if result.Error != nil {
		return nil, fmt.Errorf("ドキュメントの一覧取得失敗 %s", result.Error)
	}
	return docs, nil
}
