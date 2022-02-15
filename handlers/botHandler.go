package handlers

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// TODO вынести костантные названия кнопок в отдельный файл(Можно даже в yml)

var deleteBotMsg tgbotapi.DeleteMessageConfig

func GetMessage(telegramApiToken string) {
	bot, err := tgbotapi.NewBotAPI(telegramApiToken)
	if err != nil {
		panic(err)
	}
	bot.Debug = false
	updateConfig := tgbotapi.NewUpdate(0)

	updateConfig.Timeout = 30

	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {

		if update.CallbackQuery != nil {
			deleteBotMsg = tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
			msg, delMes := CallbackResolver(update)
			SendMessage(msg, bot)
			if delMes {
				DeleteMessage(deleteBotMsg, bot)
			}
		}

		if update.Message == nil {
			continue
		} else {
			msg = messageResolver(update)
			SendMessage(msg, bot)
		}
		//msg.ReplyToMessageID = update.Message.MessageID
	}

}

func DeleteMessage(message tgbotapi.DeleteMessageConfig, bot *tgbotapi.BotAPI) {
	if _, err := bot.Request(message); err != nil {
		fmt.Print("Error delete msg: " + err.Error())
	}
}

func SendMessage(message tgbotapi.MessageConfig, bot *tgbotapi.BotAPI) {
	_, err := bot.Send(message)
	if err != nil {
		panic("Error send msg: " + err.Error())
	}

}

//func UpdateMessage(updateMsg tgbotapi.EditMessageTextConfig, telegramApiToken string) {
//	bot, err := tgbotapi.NewBotAPI(telegramApiToken)
//	if err != nil {
//		panic(err)
//	}
//
//	if _, err := bot.Send(updateMsg); err != nil {
//		panic("Error update msg: " + err.Error())
//	}
//}
