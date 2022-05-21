package mapsActions

import (
	"errors"
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	v "github.com/spf13/viper"
	"project0/src/controllers/chatUserController"
	"project0/src/controllers/houseController"
	"project0/src/controllers/sleepUserController"
	"project0/src/controllers/userMapController"
	"project0/src/controllers/wordleController"
	"project0/src/models"
	"project0/src/repositories"
	helpers2 "project0/src/services/helpers"
)

func CheckEventUserActions(user models.User, charData []string) (msg string, buttons tg.InlineKeyboardMarkup, err error) {
	// Квесты
	switch charData[0] {

	// Дом юзера
	case v.GetString("callback_char.buy_home"):
		msg, buttons = buyHome(user)

	// Чатик
	case "chat":
		cell := models.Cell{ID: uint(helpers2.ToInt(charData[1]))}.GetCell()
		msg, buttons = chatUserController.OpenChatKeyboard(cell, user)
	case v.GetString("callback_char.join_to_chat"):
		msg, buttons = joinToChat(user, charData)

	// вордле
	case "wordle_game":
		repositories.UpdateUser(models.User{TgId: user.TgId, MenuLocation: "wordle"})
		msg, buttons = wordleController.WordleMap(user)

	// Сон
	case "sleep":
		repositories.UpdateUser(models.User{TgId: user.TgId, MenuLocation: "sleep"})
		sleepUserController.UpdateUserSleepTime(user)
		msg = sleepUserController.MsgSleepUser()
		buttons = sleepUserController.SleepButton()
	default:
		err = errors.New("not event actions")
	}

	return msg, buttons, err
}

func joinToChat(user models.User, charData []string) (msg string, buttons tg.InlineKeyboardMarkup) {
	ui := make([]models.ChatUser, 1)
	ui[0] = models.Chat{ID: uint(helpers2.ToInt(charData[1]))}.GetOrCreateChatUser(user)
	cell := models.Cell{ID: uint(helpers2.ToInt(charData[3]))}.GetCell()
	msg, buttons = chatUserController.OpenChatKeyboard(cell, user)

	helpers2.NotifyUsers(ui, v.GetString("main_info.message_user_sign_in_chat"))
	return msg, buttons
}

func buyHome(user models.User) (msg string, buttons tg.InlineKeyboardMarkup) {
	text := "Поздравляю с покупкой дома!"
	err := houseController.CreateUserHouse(user)
	if err != nil {
		switch err.Error() {
		case "user doesn't have money enough":
			text = "Не хватает деняк! Прийдется еще поднакопить :( "
		default:
			text = "Не получилось :("
		}
	}

	msg, buttons = userMapController.GetMyMap(user)
	msg = fmt.Sprintf("%s%s%s", msg, v.GetString("msg_separator"), text)
	return msg, buttons
}
