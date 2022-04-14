package repository

import (
	"fmt"
	"project0/config"
	"time"
)

type UserWords struct {
	ID       uint      `gorm:"primaryKey"`
	UserId   uint      `gorm:"embedded"`
	Word     string    `gorm:"embedded"`
	CreateAt time.Time `gorm:"autoCreateTime"`
}

func CreateUserWord(user User, word string) {
	err := config.Db.Create(&UserWords{
		Word:   word,
		UserId: user.ID,
	}).Error

	if err != nil {
		panic(err)
	}
}

func GetUserWords(user User, date time.Time) []UserWords {
	today := date.Format("2006-01-02")

	var result []UserWords
	err := config.Db.
		Where(UserWords{UserId: user.ID}).
		Where(fmt.Sprintf("create_at like '%s%s'", today, "%")).
		Find(&result).
		Error

	if err != nil {
		fmt.Println(err)
	}
	return result
}
