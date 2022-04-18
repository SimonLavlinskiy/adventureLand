package repository

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
		//Where(Instrument{ID: i.ID}).
		First(&result, Instrument{ID: i.ID}).Error

	if err != nil {
		fmt.Println("Инструмент не найден")
	}

	return result
}

func GetInstrumentsUserCanUse(user User, cell Cell) map[string]string { //todo
	instrumentsUserCanUse := map[string]string{}
	if cell.Item == nil {
		return instrumentsUserCanUse
	}
	instruments := cell.Item.Instruments

	for _, instrument := range instruments {
		if user.LeftHandId != nil && user.LeftHand.ID == instrument.Good.ID {
			instrumentsUserCanUse[user.LeftHand.View] = user.LeftHand.Type
		}
		if user.RightHandId != nil && user.RightHand.ID == instrument.Good.ID {
			instrumentsUserCanUse[user.RightHand.View] = user.RightHand.Type
		}
		if user.HeadId != nil && user.Head.ID == instrument.Good.ID {
			instrumentsUserCanUse[user.Head.View] = user.Head.Type
		}
	}
	if cell.Item.CanTake {
		instrumentsUserCanUse["👋"] = "hand"
	}
	if cell.Item.CanDestructByHand {
		instrumentsUserCanUse["🤜"] = "fist"
	}
	if cell.Item.CanStep && *cell.Type != "swap" {
		instrumentsUserCanUse["👣"] = "step"
	}

	return instrumentsUserCanUse
}
