package db

import (
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

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
