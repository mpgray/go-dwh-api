package app

import (
	"fmt"

	u "go-dwh-api/utils"
	"os"

	"github.com/go-redis/redis/v7"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB
var redisClient *redis.Client

func init() {
	postgresDB()
	redisCache()
}

func postgresDB() {
	e := godotenv.Load()
	if e != nil {
		u.Log.Warn(e.Error())
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
		u.Log.Fatal(err.Error())
	}

	db = conn
}

func redisCache() {
	//Initializing redis
	dsn := os.Getenv("REDIS_DSN")
	if len(dsn) == 0 {
		dsn = "localhost:6379"
	}
	redisClient = redis.NewClient(&redis.Options{
		Addr: dsn, //redis port
	})
	_, err := redisClient.Ping().Result()
	if err != nil {
		u.Log.Fatal(err.Error())
	}
}

// GetRedis returns the redis client
func GetRedis() *redis.Client {
	return redisClient
}

// GetDB creates the connection to our postgres database
func GetDB() *gorm.DB {
	return db
}
