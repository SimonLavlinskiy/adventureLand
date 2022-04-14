package repository

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"project0/config"
)

type UserMsg struct {
	ID       uint `gorm:"primaryKey"`
	UserTgId uint `gorm:"embedded"`
	MsgId    uint `gorm:"embedded"`
}

func UserMsgCreate(update tgbotapi.Update) {
	err := config.Db.Create(&UserMsg{
		UserTgId: uint(update.Message.Chat.ID),
		MsgId:    uint(update.Message.MessageID),
	}).Error

	if err != nil {
		panic(err)
	}
}

func IsUserMsgLatest(update tgbotapi.Update) bool {
	var lastMsg UserMsg

	config.Db.Last(&lastMsg, UserMsg{UserTgId: uint(update.Message.From.ID)})

	if lastMsg.MsgId > uint(update.Message.MessageID) {
		return false
	}

	return true
}
