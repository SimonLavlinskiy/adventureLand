package handlers

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var msg tgbotapi.MessageConfig

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

		newMessage := update.Message.Text

		switch newMessage {
		case "Карта":
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "🏔⛰🗻⬜️\U0001F7E9\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E9\U0001F7E9\U0001F7E9\U0001F7E9🚪\n⛰🗻⬜️⬜️\U0001F7E9\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E9⛪️\U0001F7E8🏪\U0001F7E9\n☃️⬜️⬜️⬜️\U0001F7E9\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E9\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E9\n⬜️⬜️⬜️🔥\U0001F7E9\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E9\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E9\n\U0001F7E9\U0001F7E9\U0001FAB5\U0001F7E9\U0001F7E9\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E9🏥\U0001F7E8🏦\U0001F7E9\n\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\n\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8🕦\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\n\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\n\U0001FAA8\U0001FAA8🐚\U0001F7E9\U0001F7E9\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E9\U0001F7E9\U0001F7E9\U0001F7E9🍄\n\U0001F7E6\U0001F7E6\U0001F7E6\U0001F7E6\U0001F7E9\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E9\U0001F7EB🍅\U0001F7EB🥔\n\U0001F7E6\U0001F7E6\U0001F7E6\U0001F7E6\U0001F7E9\U0001F7E8🐱\U0001F7E8\U0001F7E9🥕\U0001F7EB🌽\U0001F7EB\n\U0001F9CA\U0001F9CA\U0001F7E6\U0001F7E6\U0001F7E9\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E9\U0001F7EB🍎\U0001F7EB🍓\n\U0001F9CA\U0001F9CA\U0001F7E6\U0001F7E6\U0001F7E9\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E9🌳🌿🌱🌵")
		case "/start":
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Приветствую тебя, мистер "+update.Message.From.FirstName+" "+update.Message.From.LastName)
		case "/menu":
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Меню")
			//fmt.Println("1) ", tgbotapi.NewReplyKeyboard(tgbotapi.NewKeyboardButtonRow(names2buttons(mainKeyboardNames)))
			fmt.Println("2) ", mainKeyboard)
			msg.ReplyMarkup = mainKeyboard
		case "👜 Инвентарь 👜":
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Инвентарь")
			msg.ReplyMarkup = backpackKeyboard
		case "/move":
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Движение")
			msg.ReplyMarkup = moveKeyboard
		default:
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Сам ты "+update.Message.Text)
		}

		if _, err := bot.Send(msg); err != nil {
			panic(err)
		}
		//msg.ReplyToMessageID = update.Message.MessageID

	}
}
