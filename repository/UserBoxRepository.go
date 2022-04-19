package repository

import (
	"fmt"
	"project0/config"
	"time"
)

type UserBox struct {
	ID        uint      `gorm:"primaryKey"`
	UserId    uint      `gorm:"embedded"`
	BoxId     uint      `gorm:"embedded"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (userBox UserBox) IsUserGotBoxToday() bool {
	today := time.Now().Format("2006-01-02")
	var result UserBox
	err := config.Db.
		Where(fmt.Sprintf("created_at like '%s%s'", today, "%")).
		Where(userBox).
		First(&result).Error

	if err != nil {
		return false
	}

	return true
}

func (userBox UserBox) CreateUserBox() {
	err := config.Db.Create(&userBox).Error
	if err != nil {
		fmt.Println("userBox не создался")
	}
}
