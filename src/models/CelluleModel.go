package models

import (
	"fmt"
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
	NeedPay       bool `gorm:"embedded"`
	ChatId        *int `gorm:"embedded"`
	Chat          *Chat
}

func (c Cell) UpdateCell() {
	err := config.Db.
		Model(Cell{}).
		Where(&Cell{ID: c.ID}).
		Updates(&c).
		Error
	if err != nil {
		panic(err)
	}
}

func (c Cell) GetCell() (result Cell) {

	err := config.Db.
		Preload("Item").
		Preload("PrevItem").
		Preload("Teleport").
		Preload("Item.Instruments").
		Preload("Item.Instruments.Good").
		Preload("Item.Instruments.Result").
		Preload("Item.Instruments.NextStageItem").
		Preload("Chat").
		Where(c).
		First(&result).
		Error

	if err != nil {
		fmt.Println("Походу юзер вышел за границу.")
		panic(err)
	}

	return result
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
		Update("chat_id", c.ChatId).
		Error
	if err != nil {
		panic(err)
	}
}

func (c Cell) UpdateCellIfChatIsTimeout() {
	*c.ItemCount = 0

	err := config.Db.Model(Cell{}).
		Where(&Cell{ID: c.ID}).
		Update("item_id", nil).
		Update("type", "cell").
		Update("item_count", c.ItemCount).
		Update("chat_id", nil).
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

func (c Cell) IsDefaultCell() bool {
	if c.Type != nil && *c.Type == "cell" && !c.CanStep {
		return true
	}
	return false
}

func (c Cell) IsWorkbench() bool {
	if c.Type != nil && *c.Type == "workbench" && c.ItemID != nil {
		return true
	}
	return false
}

func (c Cell) IsTeleport() bool {
	if c.Type != nil && *c.Type == "teleport" && c.TeleportID != nil {
		return true
	}
	return false
}

func (c Cell) IsHome() bool {
	if c.Type != nil && *c.Type == "home" {
		return true
	}
	return false
}

func (c Cell) IsItem() bool {
	if c.Type != nil && *c.Type == "item" && c.ItemID != nil && c.ItemCount != nil && *c.ItemCount > 0 {
		return true
	}
	return false
}

func (c Cell) IsBox(user User) bool {
	if c.ItemID != nil {
		box := UserBox{UserId: user.ID, BoxId: c.Item.ID}
		if c.Type != nil && *c.Type == "item" && c.ItemID != nil && c.Item.Type == "box" && !box.IsUserGotBoxToday() {
			return true
		}
	}
	return false
}

func (c Cell) IsSwap() bool {
	if c.Type != nil && *c.Type == "swap" && c.ItemID != nil {
		return true
	}
	return false
}

func (c Cell) IsQuest() bool {
	if c.Type != nil && *c.Type == "quest" && c.ItemID != nil {
		return true
	}
	return false
}

func (c Cell) IsChat() bool {
	if c.Type != nil && *c.Type == "chat" && c.ItemID != nil && c.ChatId != nil {
		return true
	}
	return false
}

func (c Cell) IsWordleGame() bool {
	if c.Type != nil && *c.Type == "wordleGame" {
		return true
	}
	return false
}
