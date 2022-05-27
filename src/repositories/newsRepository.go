package repositories

import (
	"project0/config"
	"project0/src/models"
)

func GetNews() []models.News {
	var results []models.News

	config.Db.
		Find(&results).
		Order("id desc")

	return results
}
