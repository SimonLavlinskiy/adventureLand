package models

import (
	"project0/config"
	"time"
)

type Item struct {
	ID                  uint    `gorm:"primaryKey" json:"id"`
	Name                string  `gorm:"embedded" json:"name"`
	Description         *string `gorm:"embedded" json:"description"`
	View                string  `gorm:"embedded" json:"view"`
	Type                string  `gorm:"embedded" json:"type"`
	Cost                *int    `gorm:"embedded" json:"cost"`
	Healing             *int    `gorm:"embedded" json:"healing"`
	Damage              *int    `gorm:"embedded" json:"damage"`
	Satiety             *int    `gorm:"embedded" json:"satiety"`
	Destruction         *int    `gorm:"embedded" json:"destruction"`
	DestructionHp       *int    `gorm:"embedded" json:"destruction_hp"`
	GrowingUpTime       *int    `gorm:"embedded" json:"growing_up_time"`
	Growing             *int    `gorm:"embedded" json:"growing"`
	Breaking            *int    `gorm:"embedded" json:"breaking"`
	ItemAfterBreakingId *int    `gorm:"embedded" json:"item_after_breaking_id"`
	ItemAfterBreaking   *Item
	BreakingType        *string      `gorm:"embedded" json:"breaking_type"`
	IntervalGrowing     *int         `gorm:"embedded" json:"interval_growing"`
	CanTake             bool         `gorm:"embedded" json:"can_take"`
	CanStep             bool         `gorm:"embedded" json:"can_step"`
	CanDestructByHand   bool         `gorm:"embedded" json:"can_destruct_by_hand"`
	Instruments         []Instrument `gorm:"many2many:instrument_item;" json:"instruments"`
	DressType           *string      `gorm:"embedded" json:"dress_type"`
	IsBackpack          bool         `gorm:"embedded" json:"is_backpack"`
	IsInventory         bool         `gorm:"embedded" json:"is_inventory"`
	MaxCountUserHas     *int         `gorm:"embedded" json:"max_count_user_has"`
	CountUse            *int         `gorm:"embedded" json:"count_use"`
}

type InstrumentItem struct {
	ItemID       int        `gorm:"primaryKey" json:"item_id"`
	Item         Item       `json:"item"`
	InstrumentID int        `gorm:"primaryKey" json:"id"`
	Instrument   Instrument `json:"instrument"`
}

func (i Item) GetItemEndTime() time.Time {
	return time.Now().Add(time.Duration(*i.Growing) * time.Minute)
}

func (i Item) GetItem() Item {
	result := Item{}
	config.Db.Where(Item{ID: i.ID}).First(&result)

	return result
}
