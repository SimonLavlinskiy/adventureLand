package repository

import (
	"fmt"
	"project0/config"
	t "time"
)

type Chat struct {
	ID        uint `gorm:"primaryKey"`
	ExpiredAt t.Time
	Deleted   bool `gorm:"default:false"`
}

func CreateChat(EndTime t.Time) Chat {
	result := Chat{ExpiredAt: EndTime}
	err := config.Db.Create(&result).Error
	if err != nil {
		fmt.Println("Чатик не создался ¯ \\ _ (ツ) _ / ¯ ")
	}
	return result
}

func UpdateFiredChats() {
	var results []Chat
	config.Db.
		Where("expired_at <= ?", t.Now()).
		Where("deleted", false).
		Find(&results)

	if len(results) != 0 {
		for _, chat := range results {
			UpdateCellWithFiredChat(chat)
			//chatUser := chat.GetChatUsers()
			chat.DeleteChatUser()
			chat.DeleteChat()
		}
	}
}

func (chat Chat) DeleteChat() {
	err := config.Db.Model(Chat{}).
		Where(&Chat{ID: chat.ID}).
		Update("deleted", true).
		Error
	if err != nil {
		panic(err)
	}
}
