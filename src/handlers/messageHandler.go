package handlers

import (
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	v "github.com/spf13/viper"
	"project0/src/actions/mapsActions"
	"project0/src/controllers/learningController"
	"project0/src/controllers/sleepUserController"
	"project0/src/controllers/wordleController"
	"project0/src/models"
	"project0/src/repositories"
	"project0/src/services/helpers"
	"project0/src/services/menu"
	"strings"
)

func messageResolver(update tg.Update) (msg tg.MessageConfig) {
	user := repositories.GetOrCreateUser(update)

	fmt.Println(user.Username + " делает действие: " + msg.Text)

	if strings.Contains(user.MenuLocation, "learning") {
		msg.Text, msg.ReplyMarkup = learningController.Learning(update, user)
	}

	switch user.MenuLocation {
	case v.GetString("user_location.menu"):
		msg, msg.ReplyMarkup = menu.UserMenuLocation(update, user)
	case v.GetString("user_location.maps"):
		msg, msg.ReplyMarkup = menu.UserMapLocation(update, user)
	case v.GetString("user_location.profile"):
		msg.Text, msg.ReplyMarkup = menu.UserProfileLocation(update, user)
	case v.GetString("user_location.wordle"):
		msg.Text, msg.ReplyMarkup = wordleController.GameWordle(update, user)
	case "sleep":
		msg.Text = "\U0001F971🛌💤"
	}

	msg.ChatID = update.Message.From.ID
	msg.ParseMode = "markdown"

	return msg
}

func callBackResolver(update tg.Update) (msg tg.EditMessageTextConfig, buttons tg.EditMessageReplyMarkupConfig, newMsg bool) {
	var btns tg.InlineKeyboardMarkup

	char := update.CallbackQuery.Data
	charData := strings.Fields(update.CallbackQuery.Data)

	userTgId := helpers.GetUserTgId(update)
	user := repositories.GetUser(models.User{TgId: userTgId})

	fmt.Println(user.Username + " делает действие: " + char)

	if strings.Contains(user.MenuLocation, "learning") {
		msg.Text, btns = learningController.Learning(update, user)
	}

	// todo
	//if strings.Contains(char, "cancel") {
	// 	user = repositories.UpdateUser(models.User{TgId: user.TgId, MenuLocation: "Карта"})
	//	msg.Text, btns = userMapController.GetMyMap(user)
	//}

	switch user.MenuLocation {
	case "wordle":
		msg.Text, btns = wordleController.GameWordle(update, user)
	case "Меню":
		msg.Text, btns = menu.Menu(update, user)
	case "Профиль":
		msg.Text, btns = menu.Profile(update, user, charData)
	case "Карта":
		msg.Text, btns = mapsActions.MapsActions(user, char)
	case "sleep":
		msg.Text, btns = sleepUserController.UserSleep(user, char)
		newMsg = true
	}

	msg = tg.NewEditMessageText(update.CallbackQuery.From.ID, update.CallbackQuery.Message.MessageID, msg.Text)
	buttons = tg.NewEditMessageReplyMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, btns)
	msg.ParseMode = "markdown"

	return msg, buttons, newMsg
}
