package repository

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"project0/config"
	"strings"
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
	Shoes        *Item
	ShoesId      *int
	MenuLocation string    `gorm:"embedded"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	Deleted      bool      `gorm:"embedded"`
}

func GetOrCreateUser(update tgbotapi.Update) User {
	userId := uint(update.Message.From.ID)
	MoneyUserStart := 10

	replacer := strings.NewReplacer("_", " ", "*", " ")
	outUsername := replacer.Replace(update.Message.From.UserName)

	result := User{
		TgId:      uint(update.Message.From.ID),
		Username:  outUsername,
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
	err := config.Db.
		Preload("Head").
		Preload("RightHand").
		Preload("LeftHand").
		Preload("Body").
		Preload("Foot").
		Preload("Shoes").
		Where(user).
		First(&result).Error
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

func SetNullUserField(update tgbotapi.Update, queryFeild string) {
	var err error
	var userTgId uint

	if update.CallbackQuery != nil {
		userTgId = uint(update.CallbackQuery.From.ID)
	} else {
		userTgId = uint(update.Message.From.ID)
	}

	err = config.Db.Model(&User{}).Where(&User{TgId: userTgId}).Update(queryFeild, nil).Error

	if err != nil {
		panic(err)
	}
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
		"\n*Твое имя*: " + resUser.Username +
		"\n*Золото*: " + ToString(*resUser.Money) + "💰" +
		"\n*Аватар*: " + resUser.Avatar +
		"\n*Здоровье*: _" + ToString(int(resUser.Health)) + "_ ❤️" +
		"\n*Сытость*: _" + ToString(int(resUser.Satiety)) + "_ 😋️"

	return messageMap
}

func IsDressedItem(user User, userItem UserItem) (string, string) {
	dressItem := "Надеть"
	dressItemData := "dressGood"

	if user.HeadId != nil && userItem.ItemId == *user.HeadId ||
		user.LeftHandId != nil && userItem.ItemId == *user.LeftHandId ||
		user.RightHandId != nil && userItem.ItemId == *user.RightHandId ||
		user.BodyId != nil && userItem.ItemId == *user.BodyId ||
		user.FootId != nil && userItem.ItemId == *user.FootId ||
		user.ShoesId != nil && userItem.ItemId == *user.ShoesId {
		dressItem = "Снять"
		dressItemData = "takeOffGood"
	}

	return dressItem, dressItemData
}
