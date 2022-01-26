package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"project0/handlers"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Print("Not found .env file ")
	}

	telegramApiToken, _ := os.LookupEnv("TELEGRAM_APITOKEN")

	handlers.GetMessage(telegramApiToken)
}
