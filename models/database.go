package models

import (
	"fmt"

	u "go-dwh-api/utils"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm/logger"

	//_ "gorm.io/driver/postgres" // this is a _ 'blank' import because it is needed to load the drivers
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {

	e := godotenv.Load()
	if e != nil {
		u.Log.Warn(fmt.Sprint(e))
	}

	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")
	dbHost := os.Getenv("db_host")
	dbPort := os.Getenv("db_port")
	sslMode := os.Getenv("db_ssl_mode")

	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s port=%s sslmode=%s password=%s", dbHost, username, dbName, dbPort, sslMode, password)
	u.Log.Infof("Conntecting to DB... host=%s user=%s dbname=%s port=%s sslmode=%s password=%s", dbHost, username, dbName, dbPort, sslMode, password)

	conn, err := gorm.Open(postgres.Open(dbURI), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		u.Log.Error(fmt.Sprint(err))
	}

	db = conn
	db.AutoMigrate(&Account{}, &Contact{}, &FullName{}, &Address{}, &Phone{}, &Statement{}, &Account{})

}

// GetDB creates the connection to our postgres database
func GetDB() *gorm.DB {
	return db
}
