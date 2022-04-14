package config

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

var Db *gorm.DB
var err error

func InitMySQL() bool {

	mysqlUserName, _ := os.LookupEnv("MYSQL_USER")
	mysqlPassword, _ := os.LookupEnv("MYSQL_PASSWORD")
	mysqlHost, _ := os.LookupEnv("MYSQL_HOST")
	mysqlPort, _ := os.LookupEnv("MYSQL_PORT")
	mysqlDbName, _ := os.LookupEnv("MYSQL_DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", mysqlUserName, mysqlPassword, mysqlHost, mysqlPort, mysqlDbName)
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	//Db = Db.Debug()

	if err != nil {
		fmt.Println("DB initialization failed")
		panic(err)
		return false
	}

	return true
}
