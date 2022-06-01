package models

import "time"

type ItemCell struct {
	ID            uint `gorm:"primaryKey"`
	Cell          Cell
	CellId        int `gorm:"embedded"`
	Item          Item
	ItemId        int  `gorm:"embedded"`
	ItemCount     int  `gorm:"embedded"`
	DestructionHp *int `gorm:"embedded"`
	NextStateTime *time.Time
	LastGrowing   *time.Time
	PrevItemID    *int `gorm:"embedded"`
	PrevItem      *Item
	PrevItemCount *int `gorm:"embedded"`
}
