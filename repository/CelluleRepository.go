package repository

import (
	"fmt"
	"project0/config"
	"time"
)

type Cellule struct {
	ID            uint `gorm:"primaryKey"`
	MapsId        int  `gorm:"embedded"`
	Maps          Map
	AxisX         int     `gorm:"embedded"`
	AxisY         int     `gorm:"embedded"`
	View          string  `gorm:"embedded"`
	CanStep       bool    `gorm:"embedded"`
	Type          *string `gorm:"embedded"`
	TeleportID    *int
	Teleport      *Teleport
	ItemID        *int
	Item          *Item
	CountItem     *int `gorm:"embedded"`
	DestructionHp *int `gorm:"embedded"`
	NextStateTime *time.Time
}

func GetCellule(cellule Cellule) Cellule {
	var result Cellule

	err := config.Db.
		Preload("Item").
		Preload("Teleport").
		Preload("Item.Instruments").
		Preload("Item.Instruments.Good").
		Preload("Item.Instruments.ItemsResult").
		Preload("Item.Instruments.NextStageItem").
		Where(cellule).
		First(&result).
		Error

	if err != nil {
		fmt.Println("Походу юзер вышел за границу.")
	}

	return result
}

func UpdateCellule(cellId uint, updateCellule Cellule) {
	err := config.Db.Where(&Cellule{ID: cellId}).Updates(updateCellule).Error
	if err != nil {
		panic(err)
	}
}
