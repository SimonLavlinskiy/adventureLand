package learningController

import (
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	v "github.com/spf13/viper"
	"project0/src/actions/mapsActions"
	"project0/src/models"
	"project0/src/repositories"
	"strings"
)

func learningStep4(data string, user models.User) (text string, buttons tg.InlineKeyboardMarkup) {
	if len(data) == 0 {
		return
	}
	charData := strings.Fields(data)
	sep := v.GetString("msg_separator")

	errorInfo := "❗ Сначала тебе надо \"Надеть ✅\" все вещи на себя!"
	errorInfo2 := "❗ Пока что не надо *уничтожать 💥 предметы*!\nОни тебе еще пригодятся!"
	errorInfo3 := "❗ Пока что не надо *выкидывать 🗑 предметы*!\nОни тебе еще пригодятся!"
	info := "*Шаг 4:*\nЭто список твоих вещей, которые ты можешь использовать!\n" +
		"Используй стрелочки 🔺🔻, чтобы перемещаться по списку вещей!\n\n" +
		"- Чтобы прочитать информацию про предмет, нажми самую верхнюю кнопку «Инфо…»\n" +
		"- Чтобы использовать предмет, его надо надеть (Нажми кнопку «Надеть ✅»\n" +
		"- Чтобы скинуть предмет на карту, нужно нажать 👋🗑🗺\n" +
		"- Чтобы уничтожить предмет, нужно нажать 💥🗑💥\n" +
		"- Цифры справа от названия (30/30) показывают, сколько раз можно использовать предмет\n\n" +
		"❗️Прочитай информацию о предметах, надень их, и попробуй использовать, подойдя к другим предметам на карте"
	infoNextStep := "Используй найденные предметы, чтоб собрать ресурсы и расчистить себе путь для двери:\n" +
		"\U0001FA93 - топор для деревьев\n" +
		"⛏ - кирка для камней\n" +
		"🎩 - шляпа, чтоб прочитать информацию о предмете\n\n" +
		"Добытые ресурсы отображаются в рюкзаке"

	switch true {
	case strings.Contains(data, "cancel"):
		if userDressedAllItems(user) {
			text, buttons = mapsActions.MapsActions(user, data)
			text = fmt.Sprintf("%s%s%s", infoNextStep, sep, text)

			user.MenuLocation = "learning step5"
			repositories.UpdateUser(user)
		} else {
			text, buttons = mapsActions.MapsActions(user, "goodsMoving")
			text = fmt.Sprintf("%s%s%s%s%s", info, sep, text, sep, errorInfo)
		}
	case strings.Contains(data, "deleteItem"):
		text, buttons = mapsActions.MapsActions(user, fmt.Sprintf("goodsMoving %s", charData[2]))
		text = fmt.Sprintf("%s%s%s%s%s", info, sep, text, sep, errorInfo2)
	case strings.Contains(data, "countOfThrowOut"):
		text, buttons = mapsActions.MapsActions(user, fmt.Sprintf("goodsMoving %s", charData[2]))
		text = fmt.Sprintf("%s%s%s%s%s", info, sep, text, sep, errorInfo3)
	default:
		text, buttons = mapsActions.MapsActions(user, data)
		text = fmt.Sprintf("%s%s%s", info, sep, text)
	}

	return text, buttons
}

func userDressedAllItems(user models.User) bool {
	if user.LeftHandId != nil && user.RightHandId != nil && user.HeadId != nil {
		return true
	}
	return false
}
