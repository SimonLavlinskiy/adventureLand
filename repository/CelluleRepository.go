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
	err := config.Db.
		Where(&Cellule{ID: cellId}).
		Updates(updateCellule).
		Error
	if err != nil {
		panic(err)
	}

	if updateCellule.NextStateTime == nil {
		err := config.Db.Model(Cellule{}).
			Where(&Cellule{ID: cellId}).
			Update("next_state_time", nil).
			Error
		if err != nil {
			panic(err)
		}
	}
}

func UpdateCelluleWithNextStateTime() {
	var results []Cellule
	err := config.Db.
		Preload("Item").
		Preload("Item.Instruments").
		Preload("Item.Instruments.Good").
		Preload("Item.Instruments.ItemsResult").
		Preload("Item.Instruments.NextStageItem").
		Where("next_state_time <= ?", time.Now()).
		Find(&results).
		Error
	if err != nil || len(results) == 0 {
		fmt.Println("Ничего не найдено для обновления!")
	}

	for _, result := range results {
		for _, instrument := range result.Item.Instruments {
			if instrument.Type == "growing" {
				fmt.Println("Апдейтнулся: %?", result.ID)
				cell := updateModelCellule(result, instrument)

				UpdateCellule(result.ID,
					Cellule{
						ItemID:        cell.ItemID,
						DestructionHp: cell.DestructionHp,
						CountItem:     cell.CountItem,
						NextStateTime: cell.NextStateTime,
					})
			}
		}
	}
}

func updateModelCellule(cellule Cellule, instrument Instrument) Cellule {
	zero := 0
	if instrument.Type != "growing" && cellule.Item.DestructionHp != nil {
		cellule.DestructionHp = cellule.Item.DestructionHp
	} else if instrument.Type == "growing" && instrument.NextStageItem.DestructionHp == nil {
		cellule.DestructionHp = &zero
	} else if instrument.Type == "growing" && instrument.NextStageItem.DestructionHp != nil {
		cellule.DestructionHp = instrument.NextStageItem.DestructionHp
	}

	if *cellule.CountItem > 1 {
		*cellule.CountItem = *cellule.CountItem - 1
	} else {
		*cellule.CountItem = 0
	}

	if instrument.NextStageItem != nil {
		cellule.ItemID = instrument.NextStageItemId
	}

	if instrument.CountNextStageItem != nil {
		cellule.CountItem = instrument.CountNextStageItem
	}

	if instrument.NextStageItem != nil && instrument.NextStageItem.Growing != nil {
		*cellule.NextStateTime = time.Now().Add(time.Duration(*instrument.NextStageItem.Growing) * time.Minute)
	} else {
		cellule.NextStateTime = nil
	}

	return Cellule{
		ID:            cellule.ID,
		ItemID:        cellule.ItemID,
		CountItem:     cellule.CountItem,
		DestructionHp: cellule.DestructionHp,
		NextStateTime: cellule.NextStateTime,
	}
}
