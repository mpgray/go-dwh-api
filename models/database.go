package models

import (
	"fmt"

	u "go-hoa-api/utils"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	_ "gorm.io/driver/postgres" // this is a _ 'blank' import because it is needed to load the drivers
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {

	e := godotenv.Load()
	if e != nil {
		u.Log.Error(fmt.Sprint(e))
	}

	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")
	dbHost := os.Getenv("db_host")

	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, username, dbName, password)
	u.Log.Info(dbURI)

	conn, err := gorm.Open(postgres.Open(dbURI), &gorm.Config{})
	if err != nil {
		u.Log.Error(fmt.Sprint(err))
	}

	db = conn
	db.Debug().AutoMigrate(&Account{}, &Contact{}, &FullName{}, &Address{}, &Phone{})

}

// GetDB creates the connection to our postgres database
func GetDB() *gorm.DB {
	return db
}
