package repository

import (
	"project0/config"
)

type Receipt struct {
	ID              uint `gorm:"primaryKey"`
	Component1ID    *int `gorm:"embedded"`
	Component1      *Item
	Component1Count *int `gorm:"embedded"`
	Component2ID    *int `gorm:"embedded"`
	Component2      *Item
	Component2Count *int `gorm:"embedded"`
	Component3ID    *int `gorm:"embedded"`
	Component3      *Item
	Component3Count *int `gorm:"embedded"`
	ItemResultID    int
	ItemResult      Item
	ItemResultCount *int `gorm:"embedded"`
}

func GetReceipts() []Receipt {
	var results []Receipt
	err := config.Db.
		Preload("ItemResult").
		Preload("Component1").
		Preload("Component2").
		Preload("Component3").
		Find(&results).
		Error
	if err != nil {
		return nil
	}

	return results
}

func FindReceiptForUser(receipt Receipt) *Receipt {
	var result *Receipt
	err := config.Db.
		Preload("ItemResult").
		Preload("Component1").
		Preload("Component2").
		Preload("Component3").
		Where("component1_id", receipt.Component1ID).
		Where("component1_count", receipt.Component1Count).
		Where("component2_id", receipt.Component2ID).
		Where("component2_count", receipt.Component2Count).
		Where("component3_id", receipt.Component3ID).
		Where("component3_count", receipt.Component3Count).
		First(&result).
		Error

	if err != nil {
		return nil
	}

	return result
}
