package repository

import (
	"project0/config"
	"time"
)

type Word struct {
	ID         uint      `gorm:"primaryKey"`
	SecretWord string    `gorm:"embedded"`
	Date       time.Time `gorm:"autoCreateTime"`
}

type WordleGameProcess struct {
	ID         uint `gorm:"primaryKey"`
	User       User
	CountTries int       `gorm:"embedded"`
	Date       time.Time `gorm:"autoCreateTime"`
}

func GetActiveWord() Word {

	currentDate := time.Now()
	result := Word{}
	err := config.Db.Where(Word{Date: currentDate}).First(&result).Error

	if err != nil {
		panic(err)
	}

	return result
}
