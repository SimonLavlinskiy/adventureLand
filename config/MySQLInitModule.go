package config

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"project0/repository"
)

const driverName string = "mysql"

func InitMySQL() bool {

	mysqlUserName, _ := os.LookupEnv("MYSQL_USER")
	mysqlPassword, _ := os.LookupEnv("MYSQL_PASSWORD")
	mysqlHots, _ := os.LookupEnv("MYSQL_HOST")
	mysqlPort, _ := os.LookupEnv("MYSQL_PORT")
	mysqlDbName, _ := os.LookupEnv("MYSQL_DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", mysqlUserName, mysqlPassword, mysqlHots, mysqlPort, mysqlDbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Println("DB initialization failed")
		panic(err)
		return false
	}

	err = db.AutoMigrate(repository.User{})

	if err != nil {
		fmt.Println("Migration failed")
		panic(err)
		return false
	}
	return true
}
