package models

import (
	"fmt"
	"project0/config"
)

type Instrument struct {
	ID                 uint `gorm:"primaryKey" json:"id"`
	GoodId             *int `gorm:"embedded" json:"good_id"`
	Good               *Item
	Type               string `gorm:"embedded" json:"type"`
	ResultId           *int   `gorm:"embedded" json:"result_id"`
	Result             *Result
	NextStageItemId    *int `gorm:"embedded" json:"next_stage_item_id"`
	NextStageItem      *Item
	CountNextStageItem *int `gorm:"embedded" json:"count_next_stage_item"`
	GrowingItemId      *int `gorm:"embedded" json:"growing_item_id"`
	GrowingItem        *Item
	CountGrowingItem   *int   `gorm:"embedded" json:"count_growing_item"`
	Items              []Item `gorm:"many2many:instrument_item;"`
}

type Y struct {
	i int
	s string
}

func (i Instrument) GetInstrument() Instrument {
	var result Instrument
	err := config.Db.
		Preload("Good").
		Preload("Result").
		Preload("Result.Item").
		Preload("NextStageItem").
		Preload("GrowingItem").
		First(&result, Instrument{ID: i.ID}).Error

	if err != nil {
		fmt.Println("Инструмент не найден")
	}

	return result
}

func GetInstrumentsByItemId(itemId int) (result []Instrument) {
	err := config.Db.
		Preload("Good").
		Preload("Result").
		Preload("Result.Item").
		Preload("Result.SpecialItem").
		Preload("NextStageItem").
		Preload("GrowingItem").
		Select("i.*").
		Table("instrument_item as ii").
		Joins("left join instruments as i on i.id = ii.instrument_id").
		Where("ii.item_id = ?", itemId).
		Find(&result).
		Error

	if err != nil {
		fmt.Println("Инструмент не найден")
	}

	return result
}
