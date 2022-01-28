package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"runtime"

	"project0/handlers"

	"project0/config"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Print("Not found .env file ")
	}

	telegramApiToken, _ := os.LookupEnv("TELEGRAM_APITOKEN")

	mysqlStatus := config.InitMySQL()

	if mysqlStatus != true {
		runtime.Goexit()
	}

	handlers.GetMessage(telegramApiToken)
}
