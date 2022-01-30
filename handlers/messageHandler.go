package handlers

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	repositories "project0/repository"
	"strconv"
)

var msg tgbotapi.MessageConfig

func messageResolver(update tgbotapi.Update) tgbotapi.MessageConfig {
	res := repositories.GetOrCreateUser(update)

	newMessage := update.Message.Text

	if res.Username == "waiting" {
		res = repositories.UpdateUser(update, repositories.User{Username: newMessage})
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "*Профиль*:\n_Твое имя_ *"+res.Username+"*!\n_Аватар_:"+res.Avatar)
		msg.ParseMode = "markdown"
		msg.ReplyMarkup = profileKeyboard
	} else if res.Avatar == "waiting" {
		res = repositories.UpdateUser(update, repositories.User{Avatar: newMessage})
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "*Профиль*:\n_Твое имя_ *"+res.Username+"*!\n_Аватар_:"+res.Avatar)
		msg.ParseMode = "markdown"
		msg.ReplyMarkup = profileKeyboard
	} else {
		switch newMessage {
		case "/start":
			res := repositories.GetOrCreateUser(update)
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Приветствую тебя,  "+res.Username)
			msg.ReplyMarkup = mainKeyboard
		case "/menu", "Меню":
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Меню")

			msg.ReplyMarkup = mainKeyboard
			repositories.UpdateLocation(update, repositories.Location{Maps: "Ekaterensky"})
		case "🗺\nКарта":
			res := repositories.GetOrCreateLocation(update)
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "*Карта*: _"+res.Maps+"_; *X*: _"+strconv.FormatUint(res.AxisX, 10)+"_  *Y*: _"+strconv.FormatUint(res.AxisY, 10)+"_\n🏔⛰🗻⬜️\U0001F7E9\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E9\U0001F7E9\U0001F7E9\U0001F7E9🚪\n⛰🗻⬜️⬜️\U0001F7E9\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E9⛪️\U0001F7E8🏪\U0001F7E9\n☃️⬜️⬜️⬜️\U0001F7E9\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E9\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E9\n⬜️⬜️⬜️🔥\U0001F7E9\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E9\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E9\n\U0001F7E9\U0001F7E9\U0001FAB5\U0001F7E9\U0001F7E9\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E9🏥\U0001F7E8🏦\U0001F7E9\n\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\n\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8🕦\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\n\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\n\U0001FAA8\U0001FAA8🐚\U0001F7E9\U0001F7E9\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E9\U0001F7E9\U0001F7E9\U0001F7E9🍄\n\U0001F7E6\U0001F7E6\U0001F7E6\U0001F7E6\U0001F7E9\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E9\U0001F7EB🍅\U0001F7EB🥔\n\U0001F7E6\U0001F7E6\U0001F7E6\U0001F7E6\U0001F7E9\U0001F7E8🐱\U0001F7E8\U0001F7E9🥕\U0001F7EB🌽\U0001F7EB\n\U0001F9CA\U0001F9CA\U0001F7E6\U0001F7E6\U0001F7E9\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E9\U0001F7EB🍎\U0001F7EB🍓\n\U0001F9CA\U0001F9CA\U0001F7E6\U0001F7E6\U0001F7E9\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E9🌳🌿🌱🌵")
			msg.ReplyMarkup = moveKeyboard
		case "👤👔\nПрофиль":
			res := repositories.GetOrCreateUser(update)
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "*Профиль*:\n_Твое имя_ *"+res.Username+"*!\n_Аватар_:"+res.Avatar)
			msg.ReplyMarkup = profileKeyboard
			fmt.Print(res.Avatar, []byte(res.Avatar))
		case "📝 Изменить имя? 📝":
			res = repositories.UpdateUser(update, repositories.User{Username: "waiting"})
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "‼️ *ВНИМАНИЕ*: ‼️‼\nТы должен вписать новое имя?\n‼️‼️‼️‼️‼️‼️‼️")
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		case "👤 Изменить аватар? 👤":
			res = repositories.UpdateUser(update, repositories.User{Avatar: "waiting"})
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "‼️ *ВНИМАНИЕ*: ‼️‼\nТы должен прислать смайлик\n(_валидации пока нет_)\n‼️‼️‼️‼️‼️‼️‼️")
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		case "👜\nИнвентарь":
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Инвентарь")
			msg.ReplyMarkup = backpackKeyboard
		default:
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Сам ты "+newMessage)
		}
	}
	msg.ParseMode = "markdown"

	return msg
}
