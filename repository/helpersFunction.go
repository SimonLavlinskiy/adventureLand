package repository

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
)

func ToString(int int) string {
	return strconv.FormatInt(int64(int), 10)
}

func ToInt(string string) int {
	numInt64, _ := strconv.ParseUint(string, 10, 32)
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
