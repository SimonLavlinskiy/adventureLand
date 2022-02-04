package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"project0/repository"
	"strings"
	"time"
)

var msg tgbotapi.MessageConfig

func messageResolver(update tgbotapi.Update) tgbotapi.MessageConfig {
	resUser := repository.GetOrCreateUser(update)

	switch resUser.MenuLocation {
	case "Меню":
		msg = userMenuLocation(update, resUser)
	case "Карта":
		msg = userMapLocation(update, resUser)
	case "Профиль":
		msg = userProfileLocation(update, resUser)
	default:
		msg = userMenuLocation(update, resUser)
	}

	msg.ParseMode = "markdown"

	return msg
}

func userMenuLocation(update tgbotapi.Update, user repository.User) tgbotapi.MessageConfig {
	buttons := tgbotapi.ReplyKeyboardMarkup{}
	newMessage := update.Message.Text

	switch newMessage {
	case "🗺 Карта 🗺":
		msg.Text, buttons = repository.GetMyMap(update)
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, msg.Text)
		msg.ReplyMarkup = buttons
		repository.UpdateUser(update, repository.User{MenuLocation: "Карта"})
	case "👤 Профиль 👔":
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "*Профиль*:\n_Твое имя_ *"+user.Username+"*!\n_Аватар_:"+user.Avatar)
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, msg.Text)
		msg.ReplyMarkup = profileKeyboard
		repository.UpdateUser(update, repository.User{MenuLocation: "Профиль"})
	case "👜 Инвентарь 👜":
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Инвентарь")
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, msg.Text)
		msg.ReplyMarkup = backpackKeyboard
		repository.UpdateUser(update, repository.User{MenuLocation: "Инвентарь"})
	default:
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Меню")
		msg.ReplyMarkup = mainKeyboard
		repository.UpdateUser(update, repository.User{MenuLocation: "Меню"})
	}

	return msg
}

func userMapLocation(update tgbotapi.Update, user repository.User) tgbotapi.MessageConfig {
	newMessage := update.Message.Text
	buttons := tgbotapi.ReplyKeyboardMarkup{}
	currentTime := time.Now()
	char := strings.Fields(newMessage)

	if len(char) != 1 {
		msg = changeLocation(update, char[0])
	} else {
		switch newMessage {
		case "🔼":
			moveUp(update)
			msg.Text, buttons = repository.GetMyMap(update)
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, msg.Text)
			msg.ReplyMarkup = buttons
		case "◀️️":
			moveLeft(update)
			msg.Text, buttons = repository.GetMyMap(update)
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, msg.Text)
			msg.ReplyMarkup = buttons
		case "▶️":
			moveRight(update)
			msg.Text, buttons = repository.GetMyMap(update)
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, msg.Text)
			msg.ReplyMarkup = buttons
		case "🔽":
			moveDown(update)
			msg.Text, buttons = repository.GetMyMap(update)
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, msg.Text)
			msg.ReplyMarkup = buttons
		case "\U0001F7E6":
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Ты не похож на Jesus! 👮‍♂️")
		case "🕦":
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, currentTime.Format("\"15:04:05\"")+"\nЧасики тикают...")
		case user.Avatar:
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, repository.GetUserInfo(update))
		case "/menu", "Меню":
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Меню")
			msg.ReplyMarkup = mainKeyboard
			repository.UpdateUser(update, repository.User{MenuLocation: "Меню"})
		case "🎰":
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "💰💵🤑 Ставки на JOY CASINO дот COM! 🤑💵💰 ")
		default:
			msg.Text, buttons = repository.GetMyMap(update)
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, msg.Text)
			msg.ReplyMarkup = buttons
		}
	}

	return msg
}

func userProfileLocation(update tgbotapi.Update, user repository.User) tgbotapi.MessageConfig {
	newMessage := update.Message.Text

	if user.Username == "waiting" {
		res := repository.UpdateUser(update, repository.User{Username: newMessage})
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "*Профиль*:\n_Твое имя_ *"+res.Username+"*!\n_Аватар_:"+res.Avatar)
		msg.ReplyMarkup = profileKeyboard
	} else if user.Avatar == "waiting" {
		res := repository.UpdateUser(update, repository.User{Avatar: newMessage})
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "*Профиль*:\n_Твое имя_ *"+res.Username+"*!\n_Аватар_:"+res.Avatar)
		msg.ReplyMarkup = profileKeyboard
	} else {
		switch newMessage {
		case "📝 Изменить имя? 📝":
			repository.UpdateUser(update, repository.User{Username: "waiting"})
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "‼️ *ВНИМАНИЕ*: ‼️‼\nТы должен вписать новое имя?\n‼️‼️‼️‼️‼️‼️‼️")
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		case "👤 Изменить аватар? 👤":
			repository.UpdateUser(update, repository.User{Avatar: "waiting"})
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "‼️ *ВНИМАНИЕ*: ‼️‼\nТы должен прислать смайлик\n(_валидации пока нет_)\n‼️‼️‼️‼️‼️‼️‼️")
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		case "/menu", "Меню":
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Меню")
			msg.ReplyMarkup = mainKeyboard
			repository.UpdateUser(update, repository.User{MenuLocation: "Меню"})
		}
	}
	return msg
}

func changeLocation(update tgbotapi.Update, char string) tgbotapi.MessageConfig {
	buttons := tgbotapi.ReplyKeyboardMarkup{}

	switch char {
	case "🔼":
		moveUp(update)
	case "🔽":
		moveDown(update)
	case "◀️️":
		moveLeft(update)
	case "▶️":
		moveRight(update)
	}

	msg.Text, buttons = repository.GetMyMap(update)
	msg = tgbotapi.NewMessage(update.Message.Chat.ID, msg.Text)
	msg.ReplyMarkup = buttons

	return msg
}

func moveUp(update tgbotapi.Update) {
	res := repository.GetOrCreateMyLocation(update)
	y := *res.AxisY + 1
	repository.UpdateLocation(update, repository.Location{Map: res.Map, AxisX: res.AxisX, AxisY: &y})
}

func moveDown(update tgbotapi.Update) {
	res := repository.GetOrCreateMyLocation(update)
	y := *res.AxisY - 1
	repository.UpdateLocation(update, repository.Location{Map: res.Map, AxisX: res.AxisX, AxisY: &y})
}

func moveLeft(update tgbotapi.Update) {
	res := repository.GetOrCreateMyLocation(update)
	x := *res.AxisX - 1
	repository.UpdateLocation(update, repository.Location{Map: res.Map, AxisX: &x, AxisY: res.AxisY})
}

func moveRight(update tgbotapi.Update) {
	res := repository.GetOrCreateMyLocation(update)
	x := *res.AxisX + 1
	repository.UpdateLocation(update, repository.Location{Map: res.Map, AxisX: &x, AxisY: res.AxisY})
}
