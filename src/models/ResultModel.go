package models

import (
	"project0/config"
)

type Result struct {
	ID               uint   `gorm:"primaryKey" json:"id"`
	Type             string `gorm:"embedded" json:"type"`
	ItemId           *uint  `gorm:"embedded" json:"item_id"`
	Item             *Item
	CountItem        *uint `gorm:"embedded" json:"count_item"`
	SpecialItemId    *uint `gorm:"embedded" json:"special_item_id"`
	SpecialItem      *Item
	SpecialItemCount *uint `gorm:"embedded" json:"special_item_count"`
	Experience       *int  `gorm:"embedded" json:"experience"`
	Money            *int  `gorm:"embedded" json:"money"`
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
