package services

import (
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"os"
	"project0/repository"
)

func authTgBot(botName string) *tg.BotAPI {
	var tgApiToken string
	if botName == "AdventureLandChat" {
		tgApiToken, _ = os.LookupEnv("TELEGRAM_CHAT_APITOKEN")
	} else if botName == "AdventureLand" {
		tgApiToken, _ = os.LookupEnv("TELEGRAM_APITOKEN")
	}
	bot, err := tg.NewBotAPI(tgApiToken)
	if err != nil {
		panic(err)
	}
	bot.Debug = false

	return bot
}

func DeleteMessage(message tg.DeleteMessageConfig, bot *tg.BotAPI) {
	if _, err := bot.Request(message); err != nil {
		fmt.Print("Error delete msg: " + err.Error())
	}
}

func SendMessage(message tg.MessageConfig, bot *tg.BotAPI) tg.Message {
	resp, err := bot.Send(message)
	if err != nil {
		fmt.Printf("Error send msg: %s", err.Error())
	}
	return resp
}

//func UpdateMessage(updateMsg tg.EditMessageTextConfig, bot *tg.BotAPI) {
//	if _, err := bot.Send(updateMsg); err != nil {
//		panic(fmt.Sprintf("Error update msg: %s, %d", err.Error(), updateMsg.MessageID))
//	}
//}

func NotifyUsers(chatUsers []repository.ChatUser, message string) {
	var msg tg.MessageConfig

	bot := authTgBot("AdventureLandChat")

	for _, chatUser := range chatUsers {
		msg.Text = fmt.Sprintf("%s", message)
		msg.ChatID = int64(chatUser.User.TgId)
		msg.ReplyMarkup = tg.NewRemoveKeyboard(true)
		msg.ParseMode = "html"
		SendMessage(msg, bot)
	}
}
