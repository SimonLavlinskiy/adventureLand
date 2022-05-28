package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
	"os"
	"project0/config"
	"project0/migrations"
	"project0/src/handlers"
	"project0/templates"
	"runtime"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Print("Not found .env file ")
	}

	ViperConfiguration()

	telegramApiToken, _ := os.LookupEnv("TELEGRAM_APITOKEN")
	telegramChatApiToken, _ := os.LookupEnv("TELEGRAM_CHAT_APITOKEN")

	mysqlStatus := config.InitMySQL()
	migrations.Migrate()

	if !mysqlStatus {
		runtime.Goexit()
	}

	go handlers.RequestHandler()
	go handlers.GetMessageFromChat(telegramChatApiToken)

	handlers.GetMessage(telegramApiToken)
}

func ViperConfiguration() {
	viper.SetConfigName("templates/message")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetConfigType("yml")

	var configuration templates.Messages
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	err := viper.Unmarshal(&configuration)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}
}
