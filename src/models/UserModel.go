package models

import (
	"errors"
	"fmt"
	v "github.com/spf13/viper"
	"project0/config"
	"time"
)

type User struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	TgId       uint   `gorm:"embedded"`
	Username   string `gorm:"embedded" json:"username"`
	Avatar     string `gorm:"embedded" json:"avatar"`
	FirstName  string `gorm:"embedded"`
	LastName   string `gorm:"embedded"`
	Health     uint   `gorm:"embedded" json:"health"`
	Satiety    uint   `gorm:"embedded" json:"satiety"`
	Experience int    `gorm:"embedded" json:"experience"`
	Money      *int   `gorm:"embedded"`
	HomeId     *uint  `gorm:"embedded"`
	Clothes
	Home         *Map
	MenuLocation string    `gorm:"embedded"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"createdAt"`
	Deleted      bool      `gorm:"embedded"`
}

type Clothes struct {
	Head   *Item
	HeadId *int

	LeftHand   *Item
	LeftHandId *int

	RightHand   *Item
	RightHandId *int

	Body   *Item
	BodyId *int

	Foot   *Item
	FootId *int

	Shoes   *Item
	ShoesId *int
}

func (u User) IsDressedItem(userItem UserItem) (string, string) {
	dressItem := "Надеть ✅"
	dressItemData := v.GetString("callback_char.dress_good")

	if u.Clothes.HeadId != nil && userItem.ItemId == *u.Clothes.HeadId ||
		u.Clothes.LeftHandId != nil && userItem.ItemId == *u.Clothes.LeftHandId ||
		u.Clothes.RightHandId != nil && userItem.ItemId == *u.Clothes.RightHandId ||
		u.Clothes.BodyId != nil && userItem.ItemId == *u.Clothes.BodyId ||
		u.Clothes.FootId != nil && userItem.ItemId == *u.Clothes.FootId ||
		u.Clothes.ShoesId != nil && userItem.ItemId == *u.Clothes.ShoesId {

		dressItem = "Снять ❎"
		dressItemData = v.GetString("callback_char.take_off_good")
	}

	return dressItem, dressItemData
}

func (u User) CheckUserHasInstrument(instrument Instrument) (error, Item) {
	if instrument.Type == "hand" {
		return nil, *instrument.Good
	}
	if u.Clothes.LeftHandId != nil && *u.Clothes.LeftHandId == *instrument.GoodId {
		return nil, *u.Clothes.LeftHand
	}
	if u.Clothes.RightHandId != nil && *u.Clothes.RightHandId == *instrument.GoodId {
		return nil, *u.Clothes.RightHand
	}
	return errors.New("user dont have instrument"), Item{}
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
