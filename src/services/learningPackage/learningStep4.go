package learningPackage

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

func learningStep4(data string, user models.User) (text string, buttons tg.InlineKeyboardMarkup) {

	info := "*Шаг 3:*\nТы получил новые предметы!\nЧтобы использовать их, нажми на кнопку «Вещи 🧦»"
	infoNextStep := "Это список твоих вещей, которые ты можешь использовать!\nИспользуй стрелочки 🔺🔻, чтобы перемещаться по списку вещей!\n\n- Чтобы прочитать информацию про предмет, нажми самую верхнюю кнопку «Инфо…»\n- Чтобы использовать предмет, его надо надеть (Нажми кнопку «Надеть ✅»\n- Чтобы скинуть предмет на карту, нужно нажать 👋🗑🗺\n- Чтобы уничтожить предмет, нужно нажать 💥🗑💥\n- Цифры справа от названия (30/30) показывают, сколько раз можно использовать предмет\n\n❗️Прочитай информацию о предметах, надень их, и попробуй использовать, подойдя к другим предметам на карте"

	switch true {
	case strings.Contains(data, "goodsMoving"):
		text, buttons = mapsActions.MapsActions(user, data)
		text = fmt.Sprintf("%s%s%s", infoNextStep, v.GetString("msg_separator"), text)

		user.MenuLocation = "learning step5"
		repositories.UpdateUser(user)
	default:
		text, buttons = userMapController.GetMyMap(user)
		text = fmt.Sprintf("%s%s%s%s❗️ Нажми кнопку «Вещи 🧦»!", info, v.GetString("msg_separator"), text, v.GetString("msg_separator"))
	}

	return text, buttons
}
