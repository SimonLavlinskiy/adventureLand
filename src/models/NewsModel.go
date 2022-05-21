package models

import (
	"time"
)

type News struct {
	ID        uint      `gorm:"primaryKey"`
	Title     string    `gorm:"embedded"`
	Text      string    `gorm:"embedded"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
