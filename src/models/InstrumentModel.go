package models

import (
	"fmt"
	"project0/config"
)

type Instrument struct {
	ID                 uint `gorm:"primaryKey"`
	GoodId             *int `gorm:"embedded"`
	Good               *Item
	Type               string `gorm:"embedded"`
	ResultId           *int   `gorm:"embedded"`
	Result             *Result
	NextStageItemId    *int `gorm:"embedded"`
	NextStageItem      *Item
	CountNextStageItem *int   `gorm:"embedded"`
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
		First(&result, Instrument{ID: i.ID}).Error

	if err != nil {
		fmt.Println("Инструмент не найден")
	}

	return result
}
