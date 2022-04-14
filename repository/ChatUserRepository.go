package repository

import (
	"fmt"
	"project0/config"
)

type ChatUser struct {
	ID     uint `gorm:"primaryKey"`
	ChatID uint `gorm:"embedded"`
	Chat   Chat
	UserID uint `gorm:"embedded"`
	User   User
}

func (chat Chat) GetOrCreateChatUser(user User) ChatUser {
	result := ChatUser{
		ChatID: chat.ID,
		UserID: user.ID,
	}
	err := config.Db.
		Preload("User").
		Preload("Chat").
		Where(ChatUser{UserID: user.ID, ChatID: chat.ID}).
		FirstOrCreate(&result).Error

	if err != nil {
		fmt.Println("chat not found and not create")
	}
	return result
}

func (chat Chat) GetChatUser(user User) *ChatUser {
	result := ChatUser{}

	err := config.Db.
		Preload("User").
		Preload("Chat").
		Where(ChatUser{UserID: user.ID, ChatID: chat.ID}).
		First(&result).Error

	if err != nil {
		return nil
	}

	return &result
}

func (chat Chat) GetChatUsers() []ChatUser {
	var result []ChatUser
	config.Db.Preload("User").Where(ChatUser{ChatID: chat.ID}).Find(&result)
	return result
}

func (chat Chat) DeleteChatUser() {
	config.Db.Where("chat_id", chat.ID).Delete(ChatUser{})
}

func GetChatOfUser(user User) ChatUser {
	var result ChatUser
	config.Db.Preload("User").Where(ChatUser{UserID: user.ID}).First(&result)
	return result
}
