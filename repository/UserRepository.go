package repository

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"project0/config"
	"strconv"
	"time"
)

type User struct {
	ID           uint   `gorm:"primaryKey"`
	TgId         uint   `gorm:"embedded"`
	Username     string `gorm:"embedded"`
	Avatar       string `gorm:"embedded"`
	FirstName    string `gorm:"embedded"`
	LastName     string `gorm:"embedded"`
	Head         *Item
	HeadId       *int
	Hand         *Item
	HandId       *int
	Body         *Item
	BodyId       *int
	Foot         *Item
	FootId       *int
	MenuLocation string    `gorm:"embedded"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	Deleted      bool      `gorm:"embedded"`
}

func GetOrCreateUser(update tgbotapi.Update) User {
	result := User{
		TgId:      uint(update.Message.From.ID),
		Username:  update.Message.From.UserName,
		FirstName: update.Message.From.FirstName,
		LastName:  update.Message.From.LastName,
		Avatar:    "👤",
	}
	err := config.Db.Where(&User{TgId: uint(update.Message.From.ID)}).FirstOrCreate(&result).Error

	if err != nil {
		panic(err)
	}

	return result
}

func UpdateUser(update tgbotapi.Update, UserStruct User) User {
	var err error

	err = config.Db.Where(&User{TgId: uint(update.Message.From.ID)}).Updates(UserStruct).Error
	if err != nil {
		panic(err)
	}

	res := GetOrCreateUser(update)
	return res
}

func GetUserInfo(update tgbotapi.Update) string {
	resUser := GetOrCreateUser(update)
	resLocation := GetOrCreateMyLocation(update)

	messageMap := "*Карта*: _" + resLocation.Map + "_ *X*: _" + strconv.FormatUint(uint64(*resLocation.AxisX), 10) + "_  *Y*: _" + strconv.FormatUint(uint64(*resLocation.AxisY), 10) + "_\n_Имя_ " + resUser.Username + "\nАватар:" + resUser.Avatar

	return messageMap
}
