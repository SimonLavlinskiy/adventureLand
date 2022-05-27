package learningController

import (
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	v "github.com/spf13/viper"
	"project0/src/actions/mapsActions"
	"project0/src/controllers/mapController"
	"project0/src/models"
	"project0/src/repositories"
	"strings"
)

func learningStep5(data string, user models.User) (text string, buttons tg.InlineKeyboardMarkup) {
	if len(data) == 0 {
		return
	}
	sep := v.GetString("msg_separator")

	info := "Используй найденные предметы, чтоб собрать ресурсы и расчистить себе путь для двери:\n" +
		"\U0001FA93 - топор для деревьев\n" +
		"⛏ - кирка для камней\n" +
		"🎩 - шляпа, чтоб прочитать информацию о предмете\n\n" +
		"Добытые ресурсы отображаются в рюкзаке"

	infoNextStep := "Ой, не бойся, призраков не существует!\n" +
		"Перед тобой лежит 🎁 c 🔦 внутри, возьми его!\n" +
		"Зайди в Вещи 🧦 и Надень ✅ фонарик 🔦, и тогда сможешь развеять ночь!"

	switch true {
	case strings.Contains(data, "Меню"):
		text, buttons = mapController.GetMyMap(user)
		text = fmt.Sprintf("%s%s%s%s❗️ Пока еще рано это нажимать 🤫", info, v.GetString("msg_separator"), text, v.GetString("msg_separator"))

	case strings.Contains(data, "move 44209"):
		text, buttons = mapsActions.MapsActions(user, data)
		text = fmt.Sprintf("%s%s%s", infoNextStep, sep, text)

		user.MenuLocation = "learning step6"
		repositories.UpdateUser(user)

	default:
		text, buttons = mapsActions.MapsActions(user, data)
		text = fmt.Sprintf("%s%s%s", info, sep, text)
	}

	return text, buttons
}
