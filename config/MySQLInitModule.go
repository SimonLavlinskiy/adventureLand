package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const driverName string = "mysql"

func InitMySQL() bool {

	mysqlUserName, _ := os.LookupEnv("MYSQL_USER")
	mysqlPassword, _ := os.LookupEnv("MYSQL_PASSWORD")
	mysqlHots, _ := os.LookupEnv("MYSQL_HOST")
	mysqlPort, _ := os.LookupEnv("MYSQL_PORT")
	mysqlDbName, _ := os.LookupEnv("MYSQL_DB_NAME")

	db, err := sql.Open(driverName, fmt.Sprintf("%s:%s@/tcp(%s:%s)/%s", mysqlUserName, mysqlPassword, mysqlHots, mysqlPort, mysqlDbName))
	if err != nil {
		log.Printf("Database init error")
		log.Fatalf("error opening file: %v", err)
		return false
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return true
}
