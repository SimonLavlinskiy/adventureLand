package repository

import (
	"errors"
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
	HomeId       *uint  `gorm:"embedded"`
	Home         *Map
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
		Preload("Home").
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
		Preload("Home").
		Where(user).
		First(&result).Error
	if err != nil {
		panic(err)
	}

	return result
}

func (u User) UpdateUser() User {
	var err error

	err = config.Db.Where(&User{TgId: u.TgId}).Updates(u).Error
	if err != nil {
		panic(err)
	}

	res := GetUser(User{TgId: u.TgId})
	return res
}

func (u User) UpdateUser1() User {
	var err error
	err = config.Db.Where(&User{TgId: u.TgId}).Updates(u).Error
	if err != nil {
		panic(err)
	}

	res := GetUser(User{TgId: u.TgId})
	return res
}

func SetNullUserField(user User, queryFeild string) {
	var err error
	err = config.Db.Model(&User{}).Where(&User{TgId: user.TgId}).Update(queryFeild, nil).Error

	if err != nil {
		panic(err)
	}
}

func (u User) GetUserInfo() string {
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

func (u User) CheckUserHasInstrument(instrument Instrument) (error, Item) {
	if instrument.Type == "hand" {
		return nil, *instrument.Good
	}
	if u.LeftHandId != nil && *u.LeftHandId == *instrument.GoodId {
		return nil, *u.LeftHand
	}
	if u.RightHandId != nil && *u.RightHandId == *instrument.GoodId {
		return nil, *u.RightHand
	}
	return errors.New("User dont have instrument"), Item{}
}

func (u User) CheckUserHasLighter() (string, error) {
	if u.LeftHandId != nil && u.LeftHand.Type == "light" {

		res, err := UpdateUserInstrument(u, *u.LeftHand)
		if err != nil {
			return res, errors.New("lighter is updated")
		}

	}

	if u.RightHandId != nil && u.RightHand.Type == "light" {

		res, err := UpdateUserInstrument(u, *u.RightHand)
		if err != nil {
			return res, errors.New("lighter is updated")
		}

	}
	return "Ok", nil
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

func (u User) UserGetExperience(r Result) {
	resultExp := u.Experience + *r.Experience
	User{ID: u.ID, Experience: resultExp}.UpdateUser1()
}

func UpdateUserHand(user User, char []string) (User, UserItem) {
	userItem := UserItem{ID: ToInt(char[1])}.UserGetUserItem()

	switch char[0] {
	case v.GetString("callback_char.change_left_hand"):
		User{TgId: user.TgId, LeftHandId: &userItem.ItemId}.UpdateUser()
	case v.GetString("callback_char.change_right_hand"):
		User{TgId: user.TgId, RightHandId: &userItem.ItemId}.UpdateUser()
	}

	user = GetUser(User{TgId: user.TgId})

	return user, userItem
}
