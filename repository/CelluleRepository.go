package repository

import (
	"fmt"
	"project0/config"
)

type Cellule struct {
	ID         uint   `gorm:"primaryKey"`
	Map        string `gorm:"embedded"`
	AxisX      int    `gorm:"embedded"`
	AxisY      int    `gorm:"embedded"`
	View       string `gorm:"embedded"`
	CanStep    bool   `gorm:"embedded"`
	Type       string `gorm:"embedded"`
	TeleportId int
}

func GetCellule(cellule Cellule) Cellule {
	result := Cellule{}

	err := config.Db.Where(cellule).First(&result).Error

	if err != nil {
		fmt.Println("Походу юзер вышел за границу.")
	}

	return result
}
