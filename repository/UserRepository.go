package repository

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"project0/config"
	"project0/helpers"
	"time"
)

type User struct {
	ID           uint   `gorm:"primaryKey"`
	TgId         uint   `gorm:"embedded"`
	Username     string `gorm:"embedded"`
	Avatar       string `gorm:"embedded"`
	FirstName    string `gorm:"embedded"`
	LastName     string `gorm:"embedded"`
	Health       uint   `gorm:"embedded"`
	Satiety      uint   `gorm:"embedded"`
	Money        *int   `gorm:"embedded"`
	Head         *Item
	HeadId       *int
	LeftHand     *Item
	LeftHandId   *int
	RightHand    *Item
	RightHandId  *int
	Body         *Item
	BodyId       *int
	Foot         *Item
	FootId       *int
	MenuLocation string    `gorm:"embedded"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	Deleted      bool      `gorm:"embedded"`
}

func GetOrCreateUser(update tgbotapi.Update) User {
	userId := uint(update.Message.From.ID)
	MoneyUserStart := 0

	result := User{
		TgId:      uint(update.Message.From.ID),
		Username:  update.Message.From.UserName,
		FirstName: update.Message.From.FirstName,
		LastName:  update.Message.From.LastName,
		Avatar:    "👤",
		Satiety:   100,
		Health:    100,
		Money:     &MoneyUserStart,
	}
	err := config.Db.
		Preload("LeftHand").
		Preload("RightHand").
		Where(&User{TgId: userId}).FirstOrCreate(&result).Error

	if err != nil {
		panic(err)
	}

	return result
}

func GetUser(user User) User {
	var result User
	err := config.Db.Where(user).First(&result).Error
	if err != nil {
		panic(err)
	}
	return result
}

func UpdateUser(update tgbotapi.Update, UserStruct User) User {
	var err error
	var userTgId uint
	if update.CallbackQuery != nil {
		userTgId = uint(update.CallbackQuery.From.ID)
	} else {
		userTgId = uint(update.Message.From.ID)
	}

	err = config.Db.Where(&User{TgId: userTgId}).Updates(UserStruct).Error
	if err != nil {
		panic(err)
	}

	res := GetUser(User{TgId: userTgId})
	return res
}

func GetUserInfo(update tgbotapi.Update) string {
	var tgId uint
	if update.CallbackQuery != nil {
		tgId = uint(update.CallbackQuery.From.ID)
	} else {
		tgId = uint(update.Message.From.ID)
	}

	resUser := GetUser(User{TgId: tgId})

	messageMap := "🔅 🔆 *Профиль* 🔆 🔅\n" +
		"\n*Твое имя* " + resUser.Username +
		"\n*Золото*: " + helpers.ToString(*resUser.Money) + "💰" +
		"\n*Аватар*: " + resUser.Avatar +
		"\n*Здоровье*: _" + helpers.ToString(int(resUser.Health)) + "_ ❤️" +
		"\n*Сытость*: _" + helpers.ToString(int(resUser.Satiety)) + "_ 😋️"

	return messageMap
}
