package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"log"
	"os"
	"project0/handlers"
	"time"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Print("Not found .env file ")
	}

	telegramApiToken, _ := os.LookupEnv("TELEGRAM_APITOKEN")

	db, err := sql.Open("mysql", "user:password@/dbname")
	if err != nil {
		panic(err)
	}
	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	handlers.GetMessage(telegramApiToken)
}
