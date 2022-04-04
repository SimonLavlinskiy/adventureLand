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
	ID         uint   `gorm:"primaryKey"`
	TgId       uint   `gorm:"embedded"`
	Username   string `gorm:"embedded"`
	Avatar     string `gorm:"embedded"`
	FirstName  string `gorm:"embedded"`
	LastName   string `gorm:"embedded"`
	Health     uint   `gorm:"embedded"`
	Satiety    uint   `gorm:"embedded"`
	Experience int    `gorm:"embedded"`
	Money      *int   `gorm:"embedded"`
	Steps      uint   `gorm:"embedded"`
	HomeId     *uint  `gorm:"embedded"`
	Clothes
	Home         *Map
	MenuLocation string    `gorm:"embedded"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	Deleted      bool      `gorm:"embedded"`
	OnlineMap    *bool     `gorm:"embedded"`
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

func GetOrCreateUser(update tg.Update) User {
	userTgId := GetUserTgId(update)
	MoneyUserStart := 10
	UserOnline := true

	replacer := strings.NewReplacer("_", " ", "*", " ")
	outUsername := replacer.Replace(update.Message.From.UserName)

	result := User{
		TgId:       userTgId,
		Username:   outUsername,
		FirstName:  update.Message.From.FirstName,
		LastName:   update.Message.From.LastName,
		Avatar:     "👤",
		Satiety:    100,
		Health:     100,
		Experience: 0,
		Steps:      0,
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

func SetNullUserField(user User, queryFeild string) {
	var err error
	err = config.Db.Model(&User{}).Where(&User{TgId: user.TgId}).Update(queryFeild, nil).Error

	if err != nil {
		panic(err)
	}
}

func (u User) GetUserInfo() string {

	fmt.Println(u.Avatar)
	messageMap := fmt.Sprintf("🔅 🔆 *Профиль* 🔆 🔅\n\n"+
		"*Твое имя*: %s\n"+
		"*Аватар*: %s\n"+
		"*Золото*: %d 💰\n"+
		"*Здоровье*: _%d_ ❤️\n"+
		"*Сытость*: _%d_ 😋️\n"+
		"*Шаги*: _%d_ 👣 (_%d место_)",
		u.Username, u.Avatar, *u.Money, u.Health, u.Satiety, u.Steps, u.GetStepsPlace())

	return messageMap
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
	return errors.New("User dont have instrument"), Item{}
}

func (u User) CheckUserHasLighter() (string, error) {
	if u.Clothes.LeftHandId != nil && u.Clothes.LeftHand.Type == "light" {

		res, err := UpdateUserInstrument(u, *u.Clothes.LeftHand)
		if err != nil {
			return res, errors.New("lighter is updated")
		}

	}

	if u.Clothes.RightHandId != nil && u.Clothes.RightHand.Type == "light" {

		res, err := UpdateUserInstrument(u, *u.Clothes.RightHand)
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
	User{ID: u.ID, Experience: resultExp}.UpdateUser()
}

func UpdateUserHand(user User, char []string) (User, UserItem) {
	userItem := UserItem{ID: ToInt(char[1])}.UserGetUserItem()

	switch char[0] {
	case v.GetString("callback_char.change_left_hand"):
		clothes := &Clothes{LeftHandId: &userItem.ItemId}
		User{TgId: user.TgId, Clothes: *clothes}.UpdateUser()
	case v.GetString("callback_char.change_right_hand"):
		clothes := &Clothes{RightHandId: &userItem.ItemId}
		User{TgId: user.TgId, Clothes: *clothes}.UpdateUser()
	}

	user = GetUser(User{TgId: user.TgId})

	return user, userItem
}

func (u User) UserBuyHome(m Map) {
	*u.Money -= v.GetInt("main_info.cost_of_house")
	u.HomeId = &m.ID

	u.UpdateUser()
}

func (u User) UserStepCounter() {
	config.Db.
		Where(&User{ID: u.ID}).
		Updates(User{Steps: u.Steps + 1})
}

func (u User) GetStepsPlace() int {
	var users []User
	config.Db.
		Where(fmt.Sprintf("steps >= %d", u.Steps)).
		Find(&users)

	return len(users)
}
