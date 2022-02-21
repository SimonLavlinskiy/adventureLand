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
	LastGrowing   *time.Time
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
				UpdateCelluleAfterGrowing(result, instrument)
			}
		}
	}
}

func UpdateCelluleAfterGrowing(cellule Cellule, instrument Instrument) {

	if *cellule.CountItem > 1 {
		*cellule.CountItem = *cellule.CountItem - 1

		if cellule.Item.Growing != nil {
			*cellule.NextStateTime = time.Now().Add(time.Duration(*cellule.Item.Growing) * time.Minute)
		}
	} else {
		if instrument.NextStageItem.DestructionHp == nil {
			*cellule.DestructionHp = 0
		} else {
			cellule.DestructionHp = instrument.NextStageItem.DestructionHp
		}

		if instrument.NextStageItem != nil {
			cellule.ItemID = instrument.NextStageItemId
		}

		if instrument.CountNextStageItem != nil {
			cellule.CountItem = instrument.CountNextStageItem
		} else {
			*cellule.CountItem = 0
		}

		if instrument.NextStageItem != nil && instrument.NextStageItem.Growing != nil {
			*cellule.NextStateTime = time.Now().Add(time.Duration(*instrument.NextStageItem.Growing) * time.Minute)
		} else {
			cellule.NextStateTime = nil
		}
	}

	err := config.Db.Model(Cellule{}).
		Where(&Cellule{ID: cellule.ID}).
		Update("item_id", cellule.ItemID).
		Update("count_item", cellule.CountItem).
		Update("destruction_hp", cellule.DestructionHp).
		Update("next_state_time", cellule.NextStateTime).
		Update("last_growing", cellule.LastGrowing).
		Error
	if err != nil {
		panic(err)
	}
}

func UpdateCelluleAfterDestruction(cellule Cellule, instrument Instrument) {

	if *cellule.CountItem > 1 {
		*cellule.CountItem = *cellule.CountItem - 1

		if cellule.Item.DestructionHp != nil {
			cellule.DestructionHp = cellule.Item.DestructionHp
		}
	} else {
		if instrument.NextStageItem != nil && instrument.NextStageItem.DestructionHp != nil {
			cellule.DestructionHp = instrument.NextStageItem.DestructionHp
		}

		if instrument.NextStageItem != nil {
			cellule.ItemID = instrument.NextStageItemId
		}

		if instrument.CountNextStageItem != nil {
			cellule.CountItem = instrument.CountNextStageItem
		} else {
			*cellule.CountItem = 0
		}

		if instrument.NextStageItem != nil && instrument.NextStageItem.Growing != nil {
			*cellule.NextStateTime = time.Now().Add(time.Duration(*instrument.NextStageItem.Growing) * time.Minute)
		} else {
			cellule.NextStateTime = nil
		}
	}

	err := config.Db.Model(Cellule{}).
		Where(&Cellule{ID: cellule.ID}).
		Update("item_id", cellule.ItemID).
		Update("count_item", cellule.CountItem).
		Update("destruction_hp", cellule.DestructionHp).
		Update("next_state_time", cellule.NextStateTime).
		Error
	if err != nil {
		panic(err)
	}
}
