package repositories

import (
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"project0/config"
	"project0/src/models"
	"project0/src/services/helpers"
	"strings"
)

func GetOrCreateUser(update tg.Update) models.User {
	userTgId := helpers.GetUserTgId(update)
	MoneyUserStart := 100

	replacer := strings.NewReplacer("_", " ", "*", " ")
	outUsername := replacer.Replace(update.Message.From.UserName)

	result := models.User{
		TgId:         userTgId,
		Username:     outUsername,
		FirstName:    update.Message.From.FirstName,
		LastName:     update.Message.From.LastName,
		Avatar:       "👤",
		Satiety:      100,
		Health:       100,
		Experience:   0,
		Steps:        0,
		Money:        &MoneyUserStart,
		MenuLocation: "learning",
	}
	err := config.Db.
		Preload("Head").
		Preload("RightHand").
		Preload("LeftHand").
		Preload("Body").
		Preload("Foot").
		Preload("Shoes").
		Preload("Home").
		Where(&models.User{TgId: userTgId}).
		FirstOrCreate(&result).
		Error

	if err != nil {
		panic(err)
	}

	return result
}

func GetUser(user models.User) models.User {
	var result models.User
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

func SetNullUserField(user models.User, queryFeild string) {
	var err error
	err = config.Db.
		Model(&models.User{}).
		Where(&models.User{TgId: user.TgId}).
		Update(queryFeild, nil).
		Error

	if err != nil {
		fmt.Printf("SetNullUserField error: %s", err)
	}
}

func UpdateUser(u models.User) models.User {
	var err error

	err = config.Db.
		Where(&models.User{TgId: u.TgId}).
		Updates(&u).
		Error

	if err != nil {
		fmt.Printf("UpdateUser error: %s", err)
	}

	res := GetUser(models.User{TgId: u.TgId})

	return res
}
