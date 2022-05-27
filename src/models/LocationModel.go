package models

import (
	"project0/config"
	"time"
)

type Location struct {
	ID       uint `gorm:"primaryKey"`
	UserTgId uint
	User     User
	UserID   uint
	AxisX    *int
	AxisY    *int
	MapsId   *int `gorm:"embedded"`
	Maps     Map
	UpdateAt time.Time `gorm:"autoUpdateTime"`
}

func (u User) GetUserLocation() Location {
	result := Location{}

	err := config.Db.
		Preload("Maps").
		Where(&Location{UserID: u.ID}).
		First(result).Error
	if err != nil {
		panic(err)
	}

	return result
}
