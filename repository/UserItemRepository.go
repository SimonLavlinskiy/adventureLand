package repository

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"project0/config"
)

type UserItem struct {
	ID     int  `gorm:"primaryKey"`
	Count  *int `gorm:"embedded"`
	UserId int  `gorm:"embedded"`
	User   User
	ItemId int `gorm:"embedded"`
	Item   Item
}

func GetUserItem(update tgbotapi.Update, item Item) (UserItem, error) {
	resUser := GetOrCreateUser(update)
	userId := int(resUser.ID)
	var result UserItem
	err := config.Db.Where(UserItem{UserId: userId, ItemId: int(item.ID)}).First(&result).Error

	if err != nil {
		fmt.Println("Item not found")
	}

	return result, err
}

func GetOrCreateUserItem(update tgbotapi.Update, item Item) UserItem {
	resUser := GetOrCreateUser(update)
	userId := int(resUser.ID)
	countItem := 0
	result := UserItem{
		Count:  &countItem,
		UserId: userId,
		ItemId: int(item.ID),
	}
	err := config.Db.Where(UserItem{UserId: userId, ItemId: int(item.ID)}).FirstOrCreate(&result).Error
	if err != nil {
		panic(err)
	}

	return result
}

func GetUserItems(update tgbotapi.Update) []UserItem {
	resUser := GetOrCreateUser(update)
	userId := int(resUser.ID)
	var result []UserItem

	err := config.Db.
		Preload("Item").
		Where(UserItem{UserId: userId}).
		Where("count > 0").
		Find(&result).
		Error
	if err != nil {
		panic(err)
	}

	return result
}

func AddUserItemCount(update tgbotapi.Update, userItem UserItem, item Item) {
	resUser := GetOrCreateUser(update)
	userId := int(resUser.ID)

	sumCount := *userItem.Count + *item.Count

	err := config.Db.
		Preload("Item").
		Where(UserItem{UserId: userId, ItemId: int(item.ID)}).
		Updates(UserItem{Count: &sumCount}).
		Error
	if err != nil {
		panic(err)
	}
}
