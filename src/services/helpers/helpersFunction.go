package helpers

import (
	"encoding/json"
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
)

func ToInt(string string) int {
	numInt64, _ := strconv.ParseInt(string, 10, 32)
	return int(numInt64)
}

func GetUserTgId(update tg.Update) uint {
	var userTgId uint
	if update.CallbackQuery != nil {
		userTgId = uint(update.CallbackQuery.From.ID)
	} else if update.Message != nil {
		userTgId = uint(update.Message.From.ID)
	} else if update.MyChatMember != nil {
		userTgId = uint(update.MyChatMember.From.ID)
	}
	return userTgId
}

func CancelButton() tg.InlineKeyboardMarkup {
	return tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("Выйти", "cancel"),
		),
	)
}

func PrintJson(elem interface{}) {
	j, _ := json.Marshal(elem)
	fmt.Println(string(j))
}
