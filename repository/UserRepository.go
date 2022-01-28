package repository

import (
	"database/sql"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
	"project0/config"
	"time"
)

type User struct {
	ID         uint           `gorm:"primaryKey"`
	TgId       int64          `gorm:"embedded"`
	TgChatId   int64          `gorm:"embedded"`
	Username   string         `gorm:"embedded"`
	LocationId sql.NullString `gorm:"embedded"`
	CreatedAt  time.Time      `gorm:"autoCreateTime"`
	Deleted    bool           `gorm:"embedded"`
}

func GetOrCreateUser(update tgbotapi.Update) *gorm.DB {

	requestUser := &User{
		TgId:     update.Message.From.ID,
		TgChatId: update.Message.Chat.ID,
		Username: update.Message.From.UserName,
	}

	fmt.Println(update.Message.From.ID)
	res := config.Db.Where(User{TgId: requestUser.TgId}).First(&requestUser)
	if res.RowsAffected == 0 {
		config.Db.Create(&requestUser)
	}
	fmt.Println(res.RowsAffected)
	return res
}
