package menu

import (
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	v "github.com/spf13/viper"
	str "strings"
)

func CancelChangeNameButton(username string) tg.InlineKeyboardMarkup {
	return tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(fmt.Sprintf("или «%s» ❓", username), "cancelChangeName"),
		),
	)
}

func EmojiInlineKeyboard() tg.InlineKeyboardMarkup {
	var buttons [][]tg.InlineKeyboardButton
	var listOfAvatar []string
	listOfAvatar = str.Fields(v.GetString("message.list_of_avatar"))

	for x := 0; x < len(listOfAvatar); x = x + 8 {
		var row []tg.InlineKeyboardButton
		for i := 0; i < 8; i++ {
			sum := x + i
			if len(listOfAvatar) > sum {
				row = append(row, tg.NewInlineKeyboardButtonData(listOfAvatar[sum],
					fmt.Sprintf("%s %s", v.GetString("callback_char.change_avatar"), listOfAvatar[sum])),
				)
			}
		}
		buttons = append(buttons, row)
	}

	return tg.NewInlineKeyboardMarkup(buttons...)
}

func ProfileKeyboard(avatar string) tg.InlineKeyboardMarkup {
	return tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("📝 Изменить имя? 📝", "📝 Изменить имя? 📝"),
			tg.NewInlineKeyboardButtonData(fmt.Sprintf("%s Изменить аватар? %s", avatar, avatar), "avatarList"),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("Меню", "/menu"),
		),
	)
}

func MainKeyboard(avatar string) tg.InlineKeyboardMarkup {
	return tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("🗺 Карта 🗺", "🗺 Карта 🗺"),
			tg.NewInlineKeyboardButtonData(avatar+" Профиль 👔", avatar+" Профиль 👔"),
		),
	)
}
