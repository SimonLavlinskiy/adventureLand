package repositories

import (
	"fmt"
	"project0/config"
	"project0/src/models"
)

func GetMaps() (result []models.Map) {
	err := config.Db.Find(&result).Error

	if err != nil {
		fmt.Println("Нет карт!")
	}

	return result
}
