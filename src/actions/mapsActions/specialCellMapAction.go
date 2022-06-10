package mapsActions

import (
	"errors"
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	v "github.com/spf13/viper"
	"project0/src/controllers/mapController"
	"project0/src/controllers/newsController"
	"project0/src/models"
	"time"
)

func CheckCellEmojiAction(user models.User, charData []string) (msg string, buttons tg.InlineKeyboardMarkup, err error) {
	t := time.Now()

	// Взаимодействие с предметами на карте, у которых нет действий
	switch charData[0] {
	case v.GetString("message.emoji.water"):
		msg, buttons = useCellWithoutDoing(user, "Ты не похож на Jesus! 👮")
	case v.GetString("message.emoji.clock"):
		news := newsController.GetNewsMsg()
		msg, buttons = useCellWithoutDoing(user, fmt.Sprintf("%s\nЧасики тикают...\n\n%s", t.Format("15:04:05"), news))
	case v.GetString("message.emoji.casino"):
		msg, buttons = useCellWithoutDoing(user, "💰💵🤑 Ставки на JOY CASINO дот COM! 🤑💵💰")
	case v.GetString("message.emoji.forbidden"):
		msg, buttons = useCellWithoutDoing(user, "🚫 Сюда нельзя! 🚫")
	case v.GetString("message.emoji.shop_assistant"):
		msg, buttons = useCellWithoutDoing(user, "‍🔧 Зачем зашел за кассу? 😑")
	case v.GetString("message.emoji.wc"):
		msg, buttons = useCellWithoutDoing(user, "пись-пись 👏")
	case v.GetString("message.emoji.stop_use"):
		msg = v.GetString("errors.user_not_has_item_in_hand")
	default:
		err = errors.New("not special cell")
	}

	return msg, buttons, err
}

func useCellWithoutDoing(user models.User, text string) (msg string, buttons tg.InlineKeyboardMarkup) {
	msg, buttons = mapController.GetMyMap(user)
	msg = fmt.Sprintf("%s%s%s", msg, v.GetString("msg_separator"), text)
	return msg, buttons
}
