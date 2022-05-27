package models

import (
	"project0/config"
	"time"
)

type Item struct {
	ID                uint         `gorm:"primaryKey"`
	Name              string       `gorm:"embedded"`
	Description       *string      `gorm:"embedded"`
	View              string       `gorm:"embedded"`
	Type              string       `gorm:"embedded"`
	Cost              *int         `gorm:"embedded"`
	Healing           *int         `gorm:"embedded"`
	Damage            *int         `gorm:"embedded"`
	Satiety           *int         `gorm:"embedded"`
	Destruction       *int         `gorm:"embedded"`
	DestructionHp     *int         `gorm:"embedded"`
	GrowingUpTime     *int         `gorm:"embedded"`
	Growing           *int         `gorm:"embedded"`
	IntervalGrowing   *int         `gorm:"embedded"`
	CanTake           bool         `gorm:"embedded"`
	CanStep           bool         `gorm:"embedded"`
	CanDestructByHand bool         `gorm:"embedded"`
	Instruments       []Instrument `gorm:"many2many:instrument_item;"`
	DressType         *string      `gorm:"embedded"`
	IsBackpack        bool         `gorm:"embedded"`
	IsInventory       bool         `gorm:"embedded"`
	MaxCountUserHas   *int         `gorm:"embedded"`
	CountUse          *int         `gorm:"embedded"`
}

type InstrumentItem struct {
	ItemID       int `gorm:"primaryKey"`
	InstrumentID int `gorm:"primaryKey"`
}

func (i Item) GetItemEndTime() time.Time {
	return time.Now().Add(time.Duration(*i.Growing) * time.Minute)
}

func (i Item) GetItem() Item {
	result := Item{}
	config.Db.Where(Item{ID: i.ID}).First(&result)

	return result
}
