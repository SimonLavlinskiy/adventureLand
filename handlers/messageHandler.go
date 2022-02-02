package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"project0/repository"
	"time"
)

var msg tgbotapi.MessageConfig

func messageResolver(update tgbotapi.Update) tgbotapi.MessageConfig {
	resUser := repository.GetOrCreateUser(update)
	currentTime := time.Now()

	newMessage := update.Message.Text
	var buttons = repository.MapButtons{}

	if resUser.Username == "waiting" {
		res := repository.UpdateUser(update, repository.User{Username: newMessage})
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "*Профиль*:\n_Твое имя_ *"+res.Username+"*!\n_Аватар_:"+res.Avatar)
		msg.ParseMode = "markdown"
		msg.ReplyMarkup = profileKeyboard
	} else if resUser.Avatar == "waiting" {
		res := repository.UpdateUser(update, repository.User{Avatar: newMessage})
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "*Профиль*:\n_Твое имя_ *"+res.Username+"*!\n_Аватар_:"+res.Avatar)
		msg.ParseMode = "markdown"
		msg.ReplyMarkup = profileKeyboard
	} else {
		switch newMessage {
		case "/start":
			res := repository.GetOrCreateUser(update)
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Приветствую тебя,  "+res.Username)
			msg.ReplyMarkup = mainKeyboard
		case "/menu", "Меню":
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Меню")
			msg.ReplyMarkup = mainKeyboard
		case "🗺 Карта 🗺":
			msg.Text, buttons = repository.GetMyMap(update)
			msg.ReplyMarkup = createMoveKeyboard(buttons)
		case "👤 Профиль 👔":
			res := repository.GetOrCreateUser(update)
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "*Профиль*:\n_Твое имя_ *"+res.Username+"*!\n_Аватар_:"+res.Avatar)
			msg.ReplyMarkup = profileKeyboard
		case "📝 Изменить имя? 📝":
			repository.UpdateUser(update, repository.User{Username: "waiting"})
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "‼️ *ВНИМАНИЕ*: ‼️‼\nТы должен вписать новое имя?\n‼️‼️‼️‼️‼️‼️‼️")
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		case "👤 Изменить аватар? 👤":
			repository.UpdateUser(update, repository.User{Avatar: "waiting"})
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "‼️ *ВНИМАНИЕ*: ‼️‼\nТы должен прислать смайлик\n(_валидации пока нет_)\n‼️‼️‼️‼️‼️‼️‼️")
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		case "👜 Инвентарь 👜":
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Инвентарь")
			msg.ReplyMarkup = backpackKeyboard
		case "🔼":
			res := repository.GetOrCreateMyLocation(update)
			y := *res.AxisY + 1
			repository.UpdateLocation(update, repository.Location{Map: res.Map, AxisX: res.AxisX, AxisY: &y})
			msg.Text, buttons = repository.GetMyMap(update)
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, msg.Text)
			msg.ReplyMarkup = createMoveKeyboard(buttons)
		case "◀️":
			res := repository.GetOrCreateMyLocation(update)
			x := *res.AxisX - 1
			repository.UpdateLocation(update, repository.Location{Map: res.Map, AxisX: &x, AxisY: res.AxisY})
			msg.Text, buttons = repository.GetMyMap(update)
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, msg.Text)
			msg.ReplyMarkup = createMoveKeyboard(buttons)
		case "▶️":
			res := repository.GetOrCreateMyLocation(update)
			x := *res.AxisX + 1
			repository.UpdateLocation(update, repository.Location{Map: res.Map, AxisX: &x, AxisY: res.AxisY})
			msg.Text, buttons = repository.GetMyMap(update)
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, msg.Text)
			msg.ReplyMarkup = createMoveKeyboard(buttons)
		case "🔽":
			res := repository.GetOrCreateMyLocation(update)
			y := *res.AxisY - 1
			repository.UpdateLocation(update, repository.Location{Map: res.Map, AxisX: res.AxisX, AxisY: &y})
			msg.Text, buttons = repository.GetMyMap(update)
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, msg.Text)
			msg.ReplyMarkup = createMoveKeyboard(buttons)
		case "\U0001F7E6":
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Ты не похож на Jesus! 👮‍♂️")
		case "🕦":
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, currentTime.Format("3:4:5")+"\nЧасики тикают...")
		case resUser.Avatar:
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, repository.GetUserInfo(update))
		case "🔼🏪🚶‍♂️":
			res := repository.GetOrCreateMyLocation(update)
			y := *res.AxisY + 1
			repository.UpdateLocation(update, repository.Location{Map: "🏪 Shop", AxisX: res.AxisX, AxisY: &y})
			msg.Text, buttons = repository.GetMyMap(update)
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, msg.Text)
			msg.ReplyMarkup = createMoveKeyboard(buttons)
		case "🔼⛩🚶‍♂️":
			res := repository.GetOrCreateMyLocation(update)
			y := *res.AxisY + 1
			repository.UpdateLocation(update, repository.Location{Map: "Ekaterensky", AxisX: res.AxisX, AxisY: &y})
			msg.Text, buttons = repository.GetMyMap(update)
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, msg.Text)
			msg.ReplyMarkup = createMoveKeyboard(buttons)
		case "🔽🚪🚶‍♂️":
			res := repository.GetOrCreateMyLocation(update)
			y := *res.AxisY - 1
			repository.UpdateLocation(update, repository.Location{Map: "Main Place", AxisX: res.AxisX, AxisY: &y})
			msg.Text, buttons = repository.GetMyMap(update)
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, msg.Text)
			msg.ReplyMarkup = createMoveKeyboard(buttons)
		case "🔽⛩🚶‍♂️":
			res := repository.GetOrCreateMyLocation(update)
			y := *res.AxisY - 1
			repository.UpdateLocation(update, repository.Location{Map: "Main Place", AxisX: res.AxisX, AxisY: &y})
			msg.Text, buttons = repository.GetMyMap(update)
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, msg.Text)
			msg.ReplyMarkup = createMoveKeyboard(buttons)
		default:
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Сам ты "+newMessage)
		}
	}
	msg.ParseMode = "markdown"

	return msg
}
