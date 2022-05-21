package repositories

import (
	"project0/config"
	"project0/src/models"
	"time"
)

func GetActiveWord() (*models.Word, error) {

	currentDate := time.Now().Format("2006-01-02")
	result := models.Word{}
	err := config.Db.
		Where(models.Word{Date: currentDate}).
		First(&result).
		Error

	if err != nil {
		return nil, err
	}

	return &result, nil
}
