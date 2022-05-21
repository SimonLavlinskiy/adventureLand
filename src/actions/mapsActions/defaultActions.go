package mapsActions

import (
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	v "github.com/spf13/viper"
	"project0/src/controllers/userMapController"
	"project0/src/models"
	"project0/src/repositories"
	"project0/src/services/menu"
)

func CheckDefaultActions(user models.User, charData []string) (msg string, buttons tg.InlineKeyboardMarkup) {
	switch charData[0] {
	case "/menu", v.GetString("user_location.menu"):
		msg = "📖 Меню 📖"
		buttons = menu.MainKeyboard(user.Avatar)
		repositories.UpdateUser(models.User{TgId: user.TgId, MenuLocation: "Меню"})
	case "/map":
		msg, buttons = userMapController.GetMyMap(user)
		user = repositories.UpdateUser(models.User{TgId: user.TgId, MenuLocation: "Карта"})
	case "cancel":
		msg, buttons = userMapController.GetMyMap(user)
	default:
		msg, buttons = userMapController.GetMyMap(user)
		msg = fmt.Sprintf("%s%sХммм....🤔", msg, v.GetString("msg_separator"))
	}
	return msg, buttons
}
