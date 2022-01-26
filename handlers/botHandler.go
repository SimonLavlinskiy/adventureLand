package handlers

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var msg tgbotapi.MessageConfig

var mainKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Карта"),
		tgbotapi.NewKeyboardButton("👜 Инвентарь 👜"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Жопа"),
		tgbotapi.NewKeyboardButton("Попа"),
	),
)

var backpackKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("\U0001F9BA Шмот \U0001F9BA"),
		tgbotapi.NewKeyboardButton("🍕 Еда 🍕"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("/menu"),
	),
)

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

		if update.Message == nil {
			continue
		}

		newMessage := update.Message.Text

		if newMessage == "Карта" {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "\U0001F7E9\U0001F7E7\U0001F7E9\U0001F7E9\U0001F7E9\U0001F7E9🌲🐺\n\U0001F7E9\U0001F7E7\U0001F7E9\U0001F7E9\U0001F7E9\U0001F7E9🌲\U0001F7E9\n\U0001F7E9\U0001F7E7\U0001F7E7\U0001F7E7\U0001F7E7\U0001F7E9\U0001F7E9\U0001F7E9\n\U0001F7E9\U0001F7E7🌳🌲\U0001F7E7\U0001F7E9\U0001F7E9\U0001F7E9\n\U0001F7E9\U0001F7E7🚪🌳🐱\U0001F7E9\U0001F7E7\U0001F7E7\n\U0001F7E9\U0001FAA8🌳🌲\U0001F7E7\U0001F7E9\U0001F7E7\U0001FAA8\n\U0001FAA8\U0001F7E9\U0001F7E7\U0001F7E7\U0001F7E7\U0001F7E7\U0001F7E7\U0001FAA8\n\U0001F7E9\U0001F7E9\U0001F7E7🌳🍎\U0001F7E9\U0001F7E9\U0001F7E9")
			fmt.Println()
		} else if newMessage == "/start" {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Приветствую тебя, мистер "+update.Message.From.FirstName+" "+update.Message.From.LastName)
		} else if newMessage == "/menu" {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Меню")
			msg.ReplyMarkup = mainKeyboard
		} else if newMessage == "👜 Инвентарь 👜" {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Инвентарь")
			msg.ReplyMarkup = backpackKeyboard
		} else {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Сам ты "+update.Message.Text)
		}

		if _, err := bot.Send(msg); err != nil {
			panic(err)
		}
		//msg.ReplyToMessageID = update.Message.MessageID

	}
}
