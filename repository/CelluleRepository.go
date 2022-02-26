package repository

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
	TeleportID    *int    `gorm:"embedded"`
	Teleport      *Teleport
	ItemID        *int `gorm:"embedded"`
	Item          *Item
	ItemCount     *int `gorm:"embedded"`
	DestructionHp *int `gorm:"embedded"`
	NextStateTime *time.Time
	LastGrowing   *time.Time
	PrevItemID    *int `gorm:"embedded"`
	PrevItem      *Item
	PrevItemCount *int `gorm:"embedded"`
}

func GetCellule(cellule Cellule) Cellule {
	var result Cellule

	err := config.Db.
		Preload("Item").
		Preload("PrevItem").
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
}

func UpdateCelluleWithNextStateTime() {
	var results []Cellule
	err := config.Db.
		Preload("Item").
		Preload("PrevItem").
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
				fmt.Println("Апдейтнулся: ", result.ID)
				UpdateCelluleAfterGrowing(result, instrument)
			}
		}
	}
}

func UpdateCelluleAfterGrowing(cellule Cellule, instrument Instrument) {

	if *cellule.ItemCount > 1 && instrument.NextStageItem != nil {
		cellule.PrevItemID = cellule.ItemID
		cellule.PrevItemCount = cellule.ItemCount
		cellule = CellUpdatedNextItem(cellule, instrument)
	} else if *cellule.ItemCount > 1 && instrument.NextStageItem == nil {
		*cellule.ItemCount = *cellule.ItemCount - 1
	} else if *cellule.ItemCount <= 1 && cellule.PrevItemID != nil {
		cellule = CellUpdatePrevItem(cellule)
	} else if *cellule.ItemCount <= 1 && instrument.NextStageItem != nil {
		cellule = CellUpdatedNextItem(cellule, instrument)
	}

	cellule.LastGrowing = nil
	cellule.NextStateTime = nil

	err := config.Db.Model(Cellule{}).
		Where(&Cellule{ID: cellule.ID}).
		Update("item_id", cellule.ItemID).
		Update("item_count", cellule.ItemCount).
		Update("destruction_hp", cellule.DestructionHp).
		Update("next_state_time", cellule.NextStateTime).
		Update("last_growing", cellule.LastGrowing).
		Update("prev_item_id", cellule.PrevItemID).
		Update("prev_item_count", cellule.PrevItemCount).
		Error
	if err != nil {
		panic(err)
	}
}

func CellUpdatedNextItem(cellule Cellule, instrument Instrument) Cellule {
	cellule.ItemID = instrument.NextStageItemId

	if instrument.NextStageItem.DestructionHp == nil {
		cellule.DestructionHp = nil
	} else {
		cellule.DestructionHp = instrument.NextStageItem.DestructionHp
	}

	if instrument.CountNextStageItem == nil {
		*cellule.ItemCount = 0
	} else {
		cellule.ItemCount = instrument.CountNextStageItem
	}
	return cellule
}

func CellUpdatePrevItem(cellule Cellule) Cellule {
	cellule.ItemID = cellule.PrevItemID
	*cellule.ItemCount = *cellule.PrevItemCount - 1

	cellule.PrevItemID = nil
	cellule.PrevItemCount = nil
	cellule.NextStateTime = nil

	if cellule.PrevItem.DestructionHp != nil {
		cellule.DestructionHp = cellule.PrevItem.DestructionHp
	} else {
		cellule.DestructionHp = nil
	}
	return cellule
}

func UpdateCelluleAfterDestruction(cellule Cellule, instrument Instrument) {

	if *cellule.ItemCount > 1 && instrument.NextStageItem != nil {
		cellule.PrevItemID = cellule.ItemID
		cellule.PrevItemCount = cellule.ItemCount
		cellule = CellUpdatedNextItem(cellule, instrument)
	} else if *cellule.ItemCount > 1 && cellule.Item.DestructionHp != nil {
		*cellule.ItemCount = *cellule.ItemCount - 1
		cellule.DestructionHp = cellule.Item.DestructionHp
	} else if *cellule.ItemCount <= 1 && cellule.PrevItemID != nil {
		cellule = CellUpdatePrevItem(cellule)
	} else if *cellule.ItemCount <= 1 && instrument.NextStageItem != nil {
		cellule = CellUpdatedNextItem(cellule, instrument)
	} else {
		*cellule.ItemCount = *cellule.ItemCount - 1
		cellule.DestructionHp = nil
	}

	cellule.LastGrowing = nil
	cellule.NextStateTime = nil

	err := config.Db.Model(Cellule{}).
		Where(&Cellule{ID: cellule.ID}).
		Update("item_id", cellule.ItemID).
		Update("item_count", cellule.ItemCount).
		Update("destruction_hp", cellule.DestructionHp).
		Update("next_state_time", cellule.NextStateTime).
		Update("last_growing", cellule.LastGrowing).
		Update("prev_item_id", cellule.PrevItemID).
		Update("prev_item_count", cellule.PrevItemCount).
		Error
	if err != nil {
		panic(err)
	}
}

func UpdateCellOnPrevItem(cellule Cellule) {

	cellule = CellUpdatePrevItem(cellule)

	cellule.LastGrowing = nil
	cellule.NextStateTime = nil

	err := config.Db.Model(Cellule{}).
		Where(&Cellule{ID: cellule.ID}).
		Update("item_id", cellule.ItemID).
		Update("item_count", cellule.ItemCount).
		Update("destruction_hp", cellule.DestructionHp).
		Update("next_state_time", cellule.NextStateTime).
		Update("last_growing", cellule.LastGrowing).
		Update("prev_item_id", cellule.PrevItemID).
		Update("prev_item_count", cellule.PrevItemCount).
		Error
	if err != nil {
		panic(err)
	}
}

func UpdateCellUnderUser(update tgbotapi.Update, userItem UserItem, count int) string {
	location := GetOrCreateMyLocation(update)

	cell := GetCellule(Cellule{AxisX: *location.AxisX, AxisY: *location.AxisY, MapsId: *location.MapsId})
	if cell.ItemCount != nil && *cell.ItemCount > 0 {
		return "В этой ячейке уже есть предмет, перейди на другую ячейку..."
	}

	err := config.Db.Model(Cellule{}).
		Where(&Cellule{AxisX: *location.AxisX, AxisY: *location.AxisY, MapsId: *location.MapsId}).
		Update("item_id", userItem.ItemId).
		Update("item_count", count).
		Update("type", "item").
		Update("destruction_hp", nil).
		Update("next_state_time", nil).
		Update("last_growing", nil).
		Update("prev_item_id", nil).
		Update("prev_item_count", nil).
		Error
	if err != nil {
		panic(err)
	}

	return "Ok"

}
