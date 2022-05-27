package models

import (
	"fmt"
	"project0/config"
	"time"
)

type Word struct {
	ID         uint   `gorm:"primaryKey"`
	SecretWord string `gorm:"embedded"`
	Date       string `gorm:"embedded"`
}

type WordleGameProcess struct {
	ID         uint `gorm:"primaryKey"`
	UserId     uint `gorm:"embedded"`
	User       User
	CountTries int       `gorm:"embedded"`
	Status     string    `gorm:"embedded"`
	Date       time.Time `gorm:"autoCreateTime"`
}

func (w WordleGameProcess) UpdateWordleGameProcess(user User) {
	today := time.Now().Format("2006-01-02")

	if w.Status == "new" {
		w.CountTries++
		if w.CountTries == 6 {
			w.Status = "lose"
		}
	}

	config.Db.
		Where(&WordleGameProcess{UserId: user.ID}).
		Where(fmt.Sprintf("date like '%s%s'", today, "%")).
		Updates(WordleGameProcess{Status: w.Status, CountTries: w.CountTries})
}
