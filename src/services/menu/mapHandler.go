package menu

import (
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	v "github.com/spf13/viper"
	"project0/src/controllers/userMapController"
	"project0/src/models"
	"project0/src/repositories"
)

func UserMapLocation(update tg.Update, user models.User) (msg tg.MessageConfig, buttons tg.InlineKeyboardMarkup) {
	newMessage := update.Message.Text

	if newMessage == "/menu" || newMessage == v.GetString("user_location.menu") {
		msg.Text = "📖 Меню 📖"
		buttons = MainKeyboard(user.Avatar)
		repositories.UpdateUser(models.User{TgId: user.TgId, MenuLocation: "Меню"})
	} else if newMessage == "/map" {
		msg.Text, buttons = userMapController.GetMyMap(user)
		user = repositories.UpdateUser(models.User{TgId: user.TgId, MenuLocation: "Карта"})
	} else {
		msg.Text, buttons = userMapController.GetMyMap(user)
		msg.Text = fmt.Sprintf("%s%s🤨 Не пойму... 🧐", msg.Text, v.GetString("msg_separator"))
	}

	return msg, buttons
}
