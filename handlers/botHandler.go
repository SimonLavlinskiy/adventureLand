package handlers

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// TODO вынести костантные названия кнопок в отдельный файл(Можно даже в yml)

var mainKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("🗺 Карта 🗺"),
		tgbotapi.NewKeyboardButton("👤 Профиль 👔"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("👜 Инвентарь 👜"),
	),
)

var backpackKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("\U0001F9BA Шмот \U0001F9BA"),
		tgbotapi.NewKeyboardButton("🍕 Еда 🍕"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Меню"),
	),
)

var profileKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("📝 Изменить имя? 📝"),
		tgbotapi.NewKeyboardButton("👤 Изменить аватар? 👤"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Меню"),
	),
)

//var deleteBotMsg = tgbotapi.DeleteMessageConfig{}

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
			fmt.Println(update.CallbackQuery.Data, update.CallbackQuery.Message.MessageID)
		}

		if update.Message == nil {
			continue
		}

		//deleteBotMsg = tgbotapi.NewDeleteMessage(update.Message.Chat.ID, update.Message.MessageID-1)
		msg = messageResolver(update)

		//DeleteMessage(deleteBotMsg, telegramApiToken)
		SendMessage(msg, telegramApiToken)
		//msg.ReplyToMessageID = update.Message.MessageID

	}

}

//func DeleteMessage(message tgbotapi.DeleteMessageConfig, telegramApiToken string) {
//	bot, err := tgbotapi.NewBotAPI(telegramApiToken)
//	if err != nil {
//		panic(err)
//	}
//	if _, err := bot.Request(message); err != nil {
//		panic("Error delete msg: " + err.Error())
//	}
//}

func SendMessage(message tgbotapi.MessageConfig, telegramApiToken string) {
	bot, err := tgbotapi.NewBotAPI(telegramApiToken)
	if err != nil {
		panic(err)
	}
	if _, err := bot.Send(message); err != nil {
		panic("Error send msg: (Походу кто то спамит)" + err.Error())
	}
}

func UpdateMessage(updateMsg tgbotapi.EditMessageTextConfig, telegramApiToken string) {
	bot, err := tgbotapi.NewBotAPI(telegramApiToken)
	if err != nil {
		panic(err)
	}

	//updateMsg = tgbotapi.NewEditMessageText(update.Message.Chat.ID, 13255, "пипися")
	if _, err := bot.Send(updateMsg); err != nil {
		panic("Error update msg: " + err.Error())
	}
}
