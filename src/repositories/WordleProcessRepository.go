package repositories

import (
	"fmt"
	"project0/config"
	"project0/src/models"
	"time"
)

func GetOrCreateWordleGameProcess(user models.User) models.WordleGameProcess {
	today := time.Now().Format("2006-01-02")

	result := models.WordleGameProcess{
		UserId:     user.ID,
		CountTries: 0,
		Status:     "new",
	}

	config.Db.
		Where(&models.WordleGameProcess{UserId: user.ID}).
		Where(fmt.Sprintf("date like '%s%s'", today, "%")).
		FirstOrCreate(&result)

	return result
}

func GetWordleProcessByStatus(user models.User, status string) (result []models.WordleGameProcess) {
	config.Db.
		Where(&models.WordleGameProcess{UserId: user.ID, Status: status}).
		Find(&result)

	return result
}

func GetWordleProcessByUser(user models.User) (result []models.WordleGameProcess) {
	config.Db.
		Where(&models.WordleGameProcess{UserId: user.ID}).
		Find(&result)

	return result
}
