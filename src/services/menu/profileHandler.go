package menu

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	v "github.com/spf13/viper"
	"project0/src/controllers/userMapController"
	"project0/src/models"
	"project0/src/repositories"
	"strings"
)

func Profile(update tg.Update, user models.User, charData []string) (msg string, buttons tg.InlineKeyboardMarkup) {
	if strings.Contains(update.CallbackQuery.Data, v.GetString("callback_char.change_avatar")) {
		res := repositories.UpdateUser(models.User{TgId: user.TgId, Avatar: charData[1]})
		msg, buttons = UserProfileLocation(update, res)
	}

	switch update.CallbackQuery.Data {
	case "cancelChangeName":
		user = repositories.UpdateUser(models.User{TgId: user.TgId, Username: user.FirstName})
		msg, buttons = UserProfileLocation(update, user)
	case "📝 Изменить имя? 📝":
		repositories.UpdateUser(models.User{TgId: user.TgId, Username: "waiting"})
		msg = "‼️ *ВНИМАНИЕ*: ‼️‼\nТы должен вписать новое имя?\n‼️‼️‼️‼️‼️‼️‼️"
		buttons = CancelChangeNameButton(user.FirstName)
	case "avatarList":
		msg = "‼️ *ВНИМАНИЕ*: ‼️‼\nВыбери себе аватар..."
		buttons = EmojiInlineKeyboard()
	case "/menu", v.GetString("user_location.menu"):
		msg = "📖 Меню 📖"
		buttons = MainKeyboard(user.Avatar)
		repositories.UpdateUser(models.User{TgId: user.TgId, MenuLocation: "Меню"})
	case "/map":
		msg, buttons = userMapController.GetMyMap(user)
		user = repositories.UpdateUser(models.User{TgId: user.TgId, MenuLocation: "Карта"})
	}

	return msg, buttons
}

func UserProfileLocation(update tg.Update, user models.User) (msgText string, buttons tg.InlineKeyboardMarkup) {
	var newMessage string
	if update.Message != nil {
		newMessage = update.Message.Text
	} else {
		newMessage = update.CallbackQuery.Data
	}

	if user.Username == "waiting" {
		user = repositories.UpdateUser(models.User{TgId: user.TgId, Username: newMessage})
		msgText = user.GetUserInfo()
		buttons = ProfileKeyboard(user.Avatar)
	} else {
		switch newMessage {
		case "/map":
			msgText, buttons = userMapController.GetMyMap(user)
			user = repositories.UpdateUser(models.User{TgId: user.TgId, MenuLocation: "Карта"})
		case "/menu", v.GetString("user_location.menu"):
			msgText = "Меню"
			buttons = MainKeyboard(user.Avatar)
			repositories.UpdateUser(models.User{TgId: user.TgId, MenuLocation: "Меню"})
		default:
			msgText = user.GetUserInfo()
			buttons = ProfileKeyboard(user.Avatar)
		}
	}

	return msgText, buttons
}
