package repositories

import (
	"fmt"
	"project0/config"
	"project0/src/models"
	"time"
)

func CreateUserWord(user models.User, word string) {
	err := config.Db.Create(&models.UserWords{
		Word:   word,
		UserId: user.ID,
	}).Error

	if err != nil {
		panic(err)
	}
}

func GetUserWords(user models.User, date time.Time) []models.UserWords {
	today := date.Format("2006-01-02")

	var result []models.UserWords
	err := config.Db.
		Where(models.UserWords{UserId: user.ID}).
		Where(fmt.Sprintf("create_at like '%s%s'", today, "%")).
		Find(&result).
		Error

	if err != nil {
		fmt.Println(err)
	}
	return result
}
