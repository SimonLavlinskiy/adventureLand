package models

import (
	"fmt"
	"project0/config"
)

type Cell struct {
	ID         uint `gorm:"primaryKey"`
	MapsId     int  `gorm:"embedded"`
	Maps       Map
	AxisX      int     `gorm:"embedded"`
	AxisY      int     `gorm:"embedded"`
	View       string  `gorm:"embedded"`
	CanStep    bool    `gorm:"embedded"`
	Type       *string `gorm:"embedded"`
	TeleportID *int    `gorm:"embedded"`
	Teleport   *Teleport
	ItemCellId *uint `gorm:"embedded"`
	ItemCell   *ItemCell
	NeedPay    bool `gorm:"embedded" default:"false"`
	CanSell    bool `gorm:"embedded" default:"false"`
	ChatId     *int `gorm:"embedded"`
	Chat       *Chat
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
		Preload("Teleport").
		Preload("ItemCell.Item").
		Preload("ItemCell.ContainedItem").
		Preload("ItemCell.Item.Instruments").
		Preload("ItemCell.Item.Instruments.Good").
		Preload("ItemCell.Item.Instruments.Result").
		Preload("ItemCell.Item.Instruments.NextStageItem").
		Preload("Chat").
		Where(c).
		First(&result).
		Error

	if err != nil {
		fmt.Println(err)
	}

	return result
}

func (c Cell) UpdateCellIfChatIsTimeout() {
	*c.ItemCell.ItemCount = 0

	err := config.Db.Model(Cell{}).
		Where(&Cell{ID: c.ID}).
		Update("type", "cell").
		Update("chat_id", nil).
		Error
	if err != nil {
		panic(err)
	}

	err = config.Db.Model(ItemCell{}).
		Where(&ItemCell{ID: *c.ItemCellId}).
		Update("item_id", nil).
		Update("item_count", nil).
		Update("growing_time", nil).
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
	if c.Type != nil && *c.Type == "workbench" && c.ItemCell != nil && c.ItemCell.ItemId != nil {
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
	if c.Type != nil && *c.Type == "item" && c.ItemCell != nil && c.ItemCell.ItemId != nil && c.ItemCell.ItemCount != nil && *c.ItemCell.ItemCount > 0 {
		return true
	}
	return false
}

func (c Cell) IsBox(user User) bool {
	if c.ItemCell != nil && c.ItemCell.ItemId != nil {
		box := UserBox{UserId: user.ID, BoxId: c.ItemCell.Item.ID}
		if c.Type != nil && *c.Type == "item" && c.ItemCell.ItemId != nil && c.ItemCell.Item.Type == "box" && !box.IsUserGotBoxToday() {
			return true
		}
	}
	return false
}

func (c Cell) IsSwap() bool {
	if c.Type != nil && *c.Type == "swap" && c.ItemCell != nil && c.ItemCell.ItemId != nil {
		return true
	}
	return false
}

func (c Cell) IsQuest() bool {
	if c.Type != nil && *c.Type == "quest" && c.ItemCell != nil && c.ItemCell.ItemId != nil {
		return true
	}
	return false
}

func (c Cell) IsChat() bool {
	if c.Type != nil && *c.Type == "chat" && c.ItemCell != nil && c.ItemCell.ItemId != nil && c.ChatId != nil {
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
