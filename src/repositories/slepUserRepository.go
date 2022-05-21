package repositories

import (
	"project0/config"
	"project0/src/models"
)

func GetOrCreateUserSleep(user models.User, sleepTime models.UserSleep) {
	config.Db.
		FirstOrCreate(&sleepTime, models.UserSleep{UserId: user.ID})
}

func DeleteUserSleepTime(sleepTime models.UserSleep) {
	config.Db.
		Where(models.UserSleep{UserId: sleepTime.UserId}).
		Delete(&sleepTime)
}
