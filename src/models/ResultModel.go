package models

import (
	"project0/config"
)

type Result struct {
	ID               uint   `gorm:"primaryKey"`
	Type             string `gorm:"embedded"`
	ItemId           *uint  `gorm:"embedded"`
	Item             *Item
	CountItem        *uint `gorm:"embedded"`
	SpecialItemId    *uint `gorm:"embedded"`
	SpecialItem      *Item
	SpecialItemCount *uint `gorm:"embedded"`
	Experience       *int  `gorm:"embedded"`
}

func (r Result) GetResult() Result {
	var res Result
	config.Db.
		Preload("Item").
		Preload("SpecialItem").
		Where(Result{ID: r.ID}).
		First(&res)

	return res
}
