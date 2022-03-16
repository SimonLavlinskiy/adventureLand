package repository

import (
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"project0/config"
	"time"
)

type Cell struct {
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

func (c Cell) GetCell() Cell {
	var result Cell

	err := config.Db.
		Preload("Item").
		Preload("PrevItem").
		Preload("Teleport").
		Preload("Item.Instruments").
		Preload("Item.Instruments.Good").
		Preload("Item.Instruments.Result").
		Preload("Item.Instruments.NextStageItem").
		Where(c).
		First(&result).
		Error

	if err != nil {
		fmt.Println("Походу юзер вышел за границу.")
	}

	return result
}

func (c Cell) UpdateCell(cellId uint) {
	err := config.Db.
		Where(&Cell{ID: cellId}).
		Updates(c).
		Error
	if err != nil {
		panic(err)
	}
}

func UpdateCellWithNextStateTime() {
	var results []Cell
	err := config.Db.
		Preload("Item").
		Preload("PrevItem").
		Preload("Item.Instruments").
		Preload("Item.Instruments.Good").
		Preload("Item.Instruments.Result").
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
				result.UpdateCellAfterGrowing(instrument)
			}
		}
	}
}

func (c Cell) UpdateCellAfterGrowing(instrument Instrument) {

	if *c.ItemCount > 1 && instrument.NextStageItem != nil {
		c.PrevItemID = c.ItemID
		c.PrevItemCount = c.ItemCount
		c = c.CellUpdatedNextItem(instrument)
	} else if *c.ItemCount > 1 && instrument.NextStageItem == nil {
		*c.ItemCount = *c.ItemCount - 1
	} else if *c.ItemCount <= 1 && c.PrevItemID != nil {
		c = c.CellUpdatePrevItem()
	} else if *c.ItemCount <= 1 && instrument.NextStageItem != nil {
		c = c.CellUpdatedNextItem(instrument)
	}

	c.LastGrowing = nil
	c.NextStateTime = nil

	err := config.Db.Model(Cell{}).
		Where(&Cell{ID: c.ID}).
		Update("item_id", c.ItemID).
		Update("item_count", c.ItemCount).
		Update("destruction_hp", c.DestructionHp).
		Update("next_state_time", c.NextStateTime).
		Update("last_growing", c.LastGrowing).
		Update("prev_item_id", c.PrevItemID).
		Update("prev_item_count", c.PrevItemCount).
		Error
	if err != nil {
		panic(err)
	}
}

func (c Cell) CellUpdatedNextItem(instrument Instrument) Cell {
	c.ItemID = instrument.NextStageItemId

	if instrument.NextStageItem.DestructionHp == nil {
		c.DestructionHp = nil
	} else {
		c.DestructionHp = instrument.NextStageItem.DestructionHp
	}

	if instrument.CountNextStageItem == nil {
		*c.ItemCount = 0
	} else {
		c.ItemCount = instrument.CountNextStageItem
	}
	return c
}

func (c Cell) CellUpdatePrevItem() Cell {
	c.ItemID = c.PrevItemID
	*c.ItemCount = *c.PrevItemCount - 1

	c.PrevItemID = nil
	c.PrevItemCount = nil
	c.NextStateTime = nil

	if c.PrevItem.DestructionHp != nil {
		c.DestructionHp = c.PrevItem.DestructionHp
	} else {
		c.DestructionHp = nil
	}
	return c
}

func (c Cell) UpdateCellAfterDestruction(instrument Instrument) {

	if c.ItemCount != nil {

		if *c.ItemCount > 1 && instrument.NextStageItem != nil {
			c.PrevItemID = c.ItemID
			c.PrevItemCount = c.ItemCount
			c = c.CellUpdatedNextItem(instrument)
		} else if *c.ItemCount > 1 && c.Item.DestructionHp != nil {
			*c.ItemCount = *c.ItemCount - 1
			c.DestructionHp = c.Item.DestructionHp
		} else if *c.ItemCount <= 1 && c.PrevItemID != nil {
			c = c.CellUpdatePrevItem()
		} else if *c.ItemCount <= 1 && instrument.NextStageItem != nil {
			c = c.CellUpdatedNextItem(instrument)
		} else {
			*c.ItemCount = *c.ItemCount - 1
			c.DestructionHp = nil
		}

		c.LastGrowing = nil
		c.NextStateTime = nil

		err := config.Db.Model(Cell{}).
			Where(&Cell{ID: c.ID}).
			Update("item_id", c.ItemID).
			Update("item_count", c.ItemCount).
			Update("destruction_hp", c.DestructionHp).
			Update("next_state_time", c.NextStateTime).
			Update("last_growing", c.LastGrowing).
			Update("prev_item_id", c.PrevItemID).
			Update("prev_item_count", c.PrevItemCount).
			Error
		if err != nil {
			panic(err)
		}
	}
}

func (c Cell) UpdateCellOnPrevItem() {

	c = c.CellUpdatePrevItem()

	c.LastGrowing = nil
	c.NextStateTime = nil

	err := config.Db.Model(Cell{}).
		Where(&Cell{ID: c.ID}).
		Update("item_id", c.ItemID).
		Update("item_count", c.ItemCount).
		Update("destruction_hp", c.DestructionHp).
		Update("next_state_time", c.NextStateTime).
		Update("last_growing", c.LastGrowing).
		Update("prev_item_id", c.PrevItemID).
		Update("prev_item_count", c.PrevItemCount).
		Error
	if err != nil {
		panic(err)
	}
}

func UpdateCellUnderUser(update tg.Update, userItem UserItem, count int) string {
	location := GetOrCreateMyLocation(update)

	cell := Cell{AxisX: *location.AxisX, AxisY: *location.AxisY, MapsId: *location.MapsId}
	cell = cell.GetCell()
	if cell.ItemCount != nil && *cell.ItemCount > 0 {
		return "В этой ячейке уже есть предмет, перейди на другую ячейку..."
	}

	err := config.Db.Model(Cell{}).
		Where(&Cell{AxisX: *location.AxisX, AxisY: *location.AxisY, MapsId: *location.MapsId}).
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

func GetFullMap(id int) []Cell {
	var results []Cell

	err := config.Db.
		Preload("Item").
		Where(Cell{MapsId: id}).
		Find(&results).
		Error

	if err != nil {
		fmt.Println("Map not found!")
	}

	//j, _ := json.Marshal(results)
	//resultsJson := string(j)

	return results

}
