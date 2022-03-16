package handlers

import (
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var deleteBotMsg tg.DeleteMessageConfig

func GetMessage(telegramApiToken string) {
	bot, err := tg.NewBotAPI(telegramApiToken)
	if err != nil {
		panic(err)
	}
	bot.Debug = false
	updateConfig := tg.NewUpdate(0)

	updateConfig.Timeout = 30

	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {

		if update.CallbackQuery != nil {
			deleteBotMsg = tg.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
			msg, delMes := callBackResolver(update)
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
	}

}

func DeleteMessage(message tg.DeleteMessageConfig, bot *tg.BotAPI) {
	if _, err := bot.Request(message); err != nil {
		fmt.Print("Error delete msg: " + err.Error())
	}
}

func SendMessage(message tg.MessageConfig, bot *tg.BotAPI) {
	_, err := bot.Send(message)
	if err != nil {
		panic("Error send msg: " + err.Error())
	}
}

//func UpdateMessage(updateMsg tg.EditMessageTextConfig, telegramApiToken string) {
//	bot, err := tg.NewBotAPI(telegramApiToken)
//	if err != nil {
//		panic(err)
//	}
//
//	if _, err := bot.Send(updateMsg); err != nil {
//		panic("Error update msg: " + err.Error())
//	}
//}
