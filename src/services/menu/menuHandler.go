package menu

import (
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	v "github.com/spf13/viper"
	"project0/src/controllers/userMapController"
	"project0/src/models"
	"project0/src/repositories"
)

func Menu(update tg.Update, user models.User) (msg string, buttons tg.InlineKeyboardMarkup) {
	switch update.CallbackQuery.Data {
	case "/map":
		msg, buttons = userMapController.GetMyMap(user)
		user = repositories.UpdateUser(models.User{TgId: user.TgId, MenuLocation: "Карта"})
	case "/menu", v.GetString("user_location.menu"):
		msg = "📖 Меню 📖"
		buttons = MainKeyboard(user.Avatar)
		repositories.UpdateUser(models.User{TgId: user.TgId, MenuLocation: "Меню"})
	case "🗺 Карта 🗺":
		msg, buttons = userMapController.GetMyMap(user)
		repositories.UpdateUser(models.User{TgId: user.TgId, MenuLocation: "Карта"})
	case fmt.Sprintf("%s Профиль 👔", user.Avatar):
		msg = user.GetUserInfo()
		buttons = ProfileKeyboard(user.Avatar)
		repositories.UpdateUser(models.User{TgId: user.TgId, MenuLocation: "Профиль"})
	}

	return msg, buttons
}

func UserMenuLocation(update tg.Update, user models.User) (msg tg.MessageConfig, buttons tg.InlineKeyboardMarkup) {
	newMessage := update.Message.Text
	fmt.Println("newMsg: ", newMessage)

	switch newMessage {
	case "/userMapConfiguration":
		msg.Text, buttons = userMapController.GetMyMap(user)
		user = repositories.UpdateUser(models.User{TgId: user.TgId, MenuLocation: "Карта"})
	default:
		msg.Text = "Меню"
		buttons = MainKeyboard(user.Avatar)
		repositories.UpdateUser(models.User{TgId: user.TgId, MenuLocation: "Меню"})
	}

	return msg, buttons
}
