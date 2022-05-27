package repositories

import (
	"project0/config"
	"project0/src/models"
)

func GetReceipts() []models.Receipt {
	var results []models.Receipt
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

func FindReceiptForUser(receipt models.Receipt) *models.Receipt {
	var result *models.Receipt
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
