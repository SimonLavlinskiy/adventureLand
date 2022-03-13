package repository

import (
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	v "github.com/spf13/viper"
	"project0/config"
	"strings"
	"time"
)

type User struct {
	ID           uint   `gorm:"primaryKey"`
	TgId         uint   `gorm:"embedded"`
	TgChatId     uint   `gorm:"embedded"`
	Username     string `gorm:"embedded"`
	Avatar       string `gorm:"embedded"`
	FirstName    string `gorm:"embedded"`
	LastName     string `gorm:"embedded"`
	Health       uint   `gorm:"embedded"`
	Satiety      uint   `gorm:"embedded"`
	Experience   int    `gorm:"embedded"`
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
	OnlineMap    *bool     `gorm:"embedded"`
}

func GetOrCreateUser(update tg.Update) User {
	userTgId := GetUserTgId(update)
	MoneyUserStart := 10
	UserOnline := false

	replacer := strings.NewReplacer("_", " ", "*", " ")
	var outUsername string
	outUsername = replacer.Replace(update.Message.From.UserName)

	result := User{
		TgId:       userTgId,
		TgChatId:   uint(update.Message.Chat.ID),
		Username:   outUsername,
		FirstName:  update.Message.From.FirstName,
		LastName:   update.Message.From.LastName,
		Avatar:     "👤",
		Satiety:    100,
		Health:     100,
		Experience: 0,
		Money:      &MoneyUserStart,
		OnlineMap:  &UserOnline,
	}
	err := config.Db.
		Preload("Head").
		Preload("RightHand").
		Preload("LeftHand").
		Preload("Body").
		Preload("Foot").
		Preload("Shoes").
		Where(&User{TgId: userTgId}).
		FirstOrCreate(&result).
		Error

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

func (u User) UpdateUser(update tg.Update) User {
	var err error
	userTgId := GetUserTgId(update)
	err = config.Db.Where(&User{TgId: userTgId}).Updates(u).Error
	if err != nil {
		panic(err)
	}

	res := GetUser(User{TgId: userTgId})
	return res
}

func SetNullUserField(update tg.Update, queryFeild string) {
	var err error
	userTgId := GetUserTgId(update)
	err = config.Db.Model(&User{}).Where(&User{TgId: userTgId}).Update(queryFeild, nil).Error

	if err != nil {
		panic(err)
	}
}

func GetUserInfo(update tg.Update) string {
	userTgId := GetUserTgId(update)
	u := GetUser(User{TgId: userTgId})
	userIsOnline := "📳 Вкл"

	if !*u.OnlineMap {
		userIsOnline = "📴 Откл"
	}

	messageMap := fmt.Sprintf("🔅 🔆 *Профиль* 🔆 🔅\n\n"+
		"*Твое имя*: %s\n"+
		"*Аватар*: %s\n"+
		"*Золото*: %d 💰\n"+
		"*Здоровье*: _%d_ ❤️\n"+
		"*Сытость*: _%d_ 😋️\n"+
		"*Онлайн*: _%s_",
		u.Username, u.Avatar, *u.Money, u.Health, u.Satiety, userIsOnline)

	return messageMap
}

func (u User) IsDressedItem(userItem UserItem) (string, string) {
	dressItem := "Надеть ✅"
	dressItemData := v.GetString("callback_char.dress_good")

	if u.HeadId != nil && userItem.ItemId == *u.HeadId ||
		u.LeftHandId != nil && userItem.ItemId == *u.LeftHandId ||
		u.RightHandId != nil && userItem.ItemId == *u.RightHandId ||
		u.BodyId != nil && userItem.ItemId == *u.BodyId ||
		u.FootId != nil && userItem.ItemId == *u.FootId ||
		u.ShoesId != nil && userItem.ItemId == *u.ShoesId {

		dressItem = "Снять ❎"
		dressItemData = v.GetString("callback_char.take_off_good")
	}

	return dressItem, dressItemData
}

func (u User) CheckUserHasInstrument(instrument Instrument) (string, Item) {
	if instrument.Type == "hand" {
		return "Ok", *instrument.Good
	}
	if u.LeftHandId != nil && *u.LeftHandId == *instrument.GoodId {
		return "Ok", *u.LeftHand
	}
	if u.RightHandId != nil && *u.RightHandId == *instrument.GoodId {
		return "Ok", *u.RightHand
	}
	return "User dont have instrument", Item{}
}

func (u User) CheckUserHasLighter(update tg.Update) string {
	if u.LeftHandId != nil && u.LeftHand.Type == "light" {
		_, res := UpdateUserInstrument(update, u, *u.LeftHand)
		return res
	}
	if u.RightHandId != nil && u.RightHand.Type == "light" {
		_, res := UpdateUserInstrument(update, u, *u.RightHand)
		return res
	}
	return "Ok"
}

func (u User) GetUserQuests() []UserQuest {
	var result []UserQuest

	err := config.Db.
		Preload("Quest").
		Preload("Quest.Task").
		Preload("Quest.Result").
		Where(UserQuest{UserId: u.ID}).
		Find(&result).
		Error

	if err != nil {
		fmt.Println(fmt.Sprintf("У юзера (id = %d) нет квестов", u.ID))
	}

	return result
}
