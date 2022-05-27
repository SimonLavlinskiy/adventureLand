package repositories

import (
	"fmt"
	"project0/config"
	"project0/src/models"
	t "time"
)

func CreateChat(EndTime t.Time) (result models.Chat) {
	result.ExpiredAt = EndTime

	err := config.Db.Create(&result).Error

	if err != nil {
		fmt.Println("Чатик не создался ¯ \\ _ (ツ) _ / ¯ ")
	}

	return result
}

func GetChatOfUser(user models.User) (result models.ChatUser) {
	err := config.Db.
		Preload("User").
		Where(models.ChatUser{UserID: user.ID}).
		First(&result).
		Error

	if err != nil {
		fmt.Println("чтооо? ¯ \\ _ (ツ) _ / ¯ ")
	}

	return result
}
