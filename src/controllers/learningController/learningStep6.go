package learningController

import (
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	v "github.com/spf13/viper"
	"project0/src/actions/mapsActions"
	"project0/src/controllers/boxController"
	"project0/src/controllers/userMapController"
	"project0/src/models"
	"project0/src/repositories"
	"project0/src/services/helpers"
	"strings"
)

func learningStep6(data string, user models.User) (text string, buttons tg.InlineKeyboardMarkup) {
	if len(data) == 0 {
		return
	}

	charData := strings.Fields(data)

	sep := v.GetString("msg_separator")
	info := "Ой, не бойся, призраков не существует!\n" +
		"Перед тобой лежит 🎁, возьми его!"
	infoBox := "Зайди в Вещи 🧦 и Надень ✅ фонарик 🔦, и тогда сможешь развеять ночь!\n" +
		"_(Не забывай снимать фонарик днем, чтоб не потратить его заряд 🔋 зря!)_"
	goodInfo := "Отлично, надевай фонарик 🔦 !"
	finalInfo := "За дверью лежит прекрасный мир, который тебе еще стоит изучить!\n\n" +
		"Неподалёку от твоего места появления есть *город*!\n\n" +
		"_Не забудь сдать Квест \U0001FAA7 и получить приз за выполненное обучение!_"

	if user.MenuLocation == "learning step6.1" {
		info = infoBox
	} else if user.MenuLocation == "learning step6.2" {
		info = goodInfo
	} else if user.MenuLocation == "learning step6.3" {
		info = finalInfo
	}

	switch true {
	case strings.Contains(data, "Меню"):
		text, buttons = userMapController.GetMyMap(user)
		text = fmt.Sprintf("%s%s%s%s❗️ Пока еще рано это нажимать 🤫", info, v.GetString("msg_separator"), text, v.GetString("msg_separator"))

	case strings.Contains(data, "move 44316"):
		text, buttons = mapsActions.MapsActions(user, data)
		text = fmt.Sprintf("%s%s%s", finalInfo, sep, text)

		user.MenuLocation = "Карта"
		repositories.UpdateUser(user)

	case strings.Contains(data, "box"):
		cell := models.Cell{ID: uint(helpers.ToInt(charData[1]))}.GetCell()
		text, buttons = boxController.UserGetBox(user, cell)
		text = fmt.Sprintf("%s%s%s", infoBox, v.GetString("msg_separator"), text)

		user.MenuLocation = "learning step6.1"
		repositories.UpdateUser(user)

	case strings.Contains(data, "goodsMoving") && user.MenuLocation == "learning step6.1":
		text, buttons = mapsActions.MapsActions(user, data)
		text = fmt.Sprintf("%s%s%s", goodInfo, v.GetString("msg_separator"), text)

		user.MenuLocation = "learning step6.2"
		repositories.UpdateUser(user)

	case strings.Contains(data, "cancel") && user.MenuLocation == "learning step6.2":
		text, buttons = mapsActions.MapsActions(user, data)
		text = fmt.Sprintf("%s%s%s", finalInfo, sep, text)

		user.MenuLocation = "learning step6.3"
		repositories.UpdateUser(user)

	case strings.Contains(data, "dressGood"):
		text, buttons = mapsActions.MapsActions(user, data)

	default:
		text, buttons = mapsActions.MapsActions(user, data)
		text = fmt.Sprintf("%s%s%s", info, sep, text)
	}

	return text, buttons
}
