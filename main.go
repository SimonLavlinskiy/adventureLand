package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"project0/config"
	"project0/handlers"
	"project0/migrations"
	"runtime"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Print("Not found .env file ")
	}

	telegramApiToken, _ := os.LookupEnv("TELEGRAM_APITOKEN")

	mysqlStatus := config.InitMySQL()
	migrations.Migrate()

	if mysqlStatus != true {
		runtime.Goexit()
	}

	handlers.GetMessage(telegramApiToken)
}
