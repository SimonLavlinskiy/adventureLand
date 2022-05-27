package models

import (
	"time"
)

type UserWords struct {
	ID       uint      `gorm:"primaryKey"`
	UserId   uint      `gorm:"embedded"`
	Word     string    `gorm:"embedded"`
	CreateAt time.Time `gorm:"autoCreateTime"`
}
