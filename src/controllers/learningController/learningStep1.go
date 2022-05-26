package learningController

import (
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	v "github.com/spf13/viper"
	"project0/src/actions/mapsActions"
	"project0/src/controllers/userMapController"
	"project0/src/models"
	"project0/src/repositories"
	"strings"
)

func learningStep1(data string, user models.User) (text string, buttons tg.InlineKeyboardMarkup) {
	info := "*Шаг 1:*\nВидишь снизу кнопки-стрелочки (◀️ 🔼 ▶️ 🔽)? Они позволяют тебе ходить!\nПопробуй пройтись по карте, а как освоишься, бери квест \U0001FAA7 на обучение и заходи в дверь 🚪"
	infoNextStep := "*Шаг 2:*\nЗдесь ты научишься добывать ресурсы!\nИсследуй местность, не бойся нажимать на кнопки и не забудь взять подарочки 🎁 сверху, там ты получишь инструменты!"

	switch true {
	case strings.Contains(data, "goodsMoving"), strings.Contains(data, "Меню"), strings.Contains(data, "category"):
		text, buttons = userMapController.GetMyMap(user)
		text = fmt.Sprintf("%s%s%s%s❗️ Пока еще рано это нажимать 🤫", info, v.GetString("msg_separator"), text, v.GetString("msg_separator"))
	case strings.Contains(data, "move 22"):
		user.MenuLocation = "learning step2"
		repositories.UpdateUser(user)

		text, buttons = mapsActions.MapsActions(user, data)
		text = fmt.Sprintf("%s%s%s", infoNextStep, v.GetString("msg_separator"), text)
	default:
		text, buttons = mapsActions.MapsActions(user, data)
		text = fmt.Sprintf("%s%s%s", info, v.GetString("msg_separator"), text)
	}

	return text, buttons
}
