package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var msg tgbotapi.MessageConfig

func messageResolver(update tgbotapi.Update) tgbotapi.MessageConfig {

	newMessage := update.Message.Text

	switch newMessage {
	case "Карта":
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "🏔⛰🗻⬜️\U0001F7E9\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E9\U0001F7E9\U0001F7E9\U0001F7E9🚪\n⛰🗻⬜️⬜️\U0001F7E9\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E9⛪️\U0001F7E8🏪\U0001F7E9\n☃️⬜️⬜️⬜️\U0001F7E9\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E9\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E9\n⬜️⬜️⬜️🔥\U0001F7E9\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E9\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E9\n\U0001F7E9\U0001F7E9\U0001FAB5\U0001F7E9\U0001F7E9\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E9🏥\U0001F7E8🏦\U0001F7E9\n\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\n\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8🕦\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\n\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\n\U0001FAA8\U0001FAA8🐚\U0001F7E9\U0001F7E9\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E9\U0001F7E9\U0001F7E9\U0001F7E9🍄\n\U0001F7E6\U0001F7E6\U0001F7E6\U0001F7E6\U0001F7E9\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E9\U0001F7EB🍅\U0001F7EB🥔\n\U0001F7E6\U0001F7E6\U0001F7E6\U0001F7E6\U0001F7E9\U0001F7E8🐱\U0001F7E8\U0001F7E9🥕\U0001F7EB🌽\U0001F7EB\n\U0001F9CA\U0001F9CA\U0001F7E6\U0001F7E6\U0001F7E9\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E9\U0001F7EB🍎\U0001F7EB🍓\n\U0001F9CA\U0001F9CA\U0001F7E6\U0001F7E6\U0001F7E9\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E9🌳🌿🌱🌵")
		msg.ReplyMarkup = moveKeyboard
	case "/start":
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Приветствую тебя, мистер "+update.Message.From.FirstName+" "+update.Message.From.LastName)
	case "/menu":
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Меню")
		msg.ReplyMarkup = mainKeyboard
	case "👜 Инвентарь 👜":
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Инвентарь")
		msg.ReplyMarkup = backpackKeyboard
	default:
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Сам ты "+newMessage)
	}

	return msg
}
