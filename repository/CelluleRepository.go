package repository

import (
	"fmt"
	"project0/config"
)

type Cellule struct {
	ID         uint    `gorm:"primaryKey"`
	Map        string  `gorm:"embedded"`
	AxisX      int     `gorm:"embedded"`
	AxisY      int     `gorm:"embedded"`
	View       string  `gorm:"embedded"`
	CanStep    bool    `gorm:"embedded"`
	Type       *string `gorm:"embedded"`
	TeleportID *int
	Teleport   *Teleport
	ItemID     *int
	Item       *Item
}

func GetCellule(cellule Cellule) Cellule {
	var result Cellule

	err := config.Db.
		Preload("Item").
		Preload("Teleport").
		Where(cellule).
		First(&result).
		Error

	//j, _ := json.Marshal(result)
	//fmt.Println(string(j))

	if err != nil {
		fmt.Println("Походу юзер вышел за границу.")
	}

	return result
}