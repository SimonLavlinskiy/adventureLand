package repository

import (
	"encoding/json"
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

func UpdateCelluleWhenUserUseIt(cellule Cellule, instrument Instrument) {
	now := time.Now()
	nextTimeStateItem := now
	var destructionHp int

	if instrument.NextStageTimeMin != nil {
		nextTimeStateItem = now.Add(time.Duration(*instrument.NextStageTimeMin) * time.Minute)
	}
	j, _ := json.Marshal(instrument)
	fmt.Println(string(j))
	if instrument.NextStageItemId != nil && instrument.ItemsResult.DestructionHp != nil {
		destructionHp = *cellule.Item.DestructionHp
	}

	var cellUpdate = Cellule{
		ItemID:        instrument.NextStageItemId,
		CountItem:     instrument.CountNextStageItem,
		DestructionHp: &destructionHp,
		NextStateTime: &nextTimeStateItem,
	}
	UpdateCellule(cellule.ID, cellUpdate)
}
