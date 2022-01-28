package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// TODO вынести костантные названия кнопок в отдельный файл(Можно даже в yml)

var mainKeyboardNames = []string{
	"Карта", "👜 Инвентарь 👜",
}

var menuButtons = []string{
	"/menu",
}

var backpackKeyboardNames = [][]string{
	{"\U0001F9BA Шмот \U0001F9BA", "\"🍕 Еда 🍕\""},
}

func names2buttons(names []string) []tgbotapi.KeyboardButton {
	var row []tgbotapi.KeyboardButton
	for _, l := range names {
		row = append(row, tgbotapi.NewKeyboardButton(l))
	}
	return row
}

var mainKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Карта"),
		tgbotapi.NewKeyboardButton("👜 Инвентарь 👜"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("/menu"),
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

var moveKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("⬛"),
		tgbotapi.NewKeyboardButton("🔼"),
		tgbotapi.NewKeyboardButton("⬛"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("◀️"),
		tgbotapi.NewKeyboardButton("️⏺"),
		tgbotapi.NewKeyboardButton("▶️"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("⬛"),
		tgbotapi.NewKeyboardButton("🔽"),
		tgbotapi.NewKeyboardButton("/menu"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("\U0001F9BA Шмот \U0001F9BA"),
		tgbotapi.NewKeyboardButton("🍕 Еда 🍕"),
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

		msg := messageResolver(update)

		if _, err := bot.Send(msg); err != nil {
			panic(err)
		}
		//msg.ReplyToMessageID = update.Message.MessageID

	}
}
