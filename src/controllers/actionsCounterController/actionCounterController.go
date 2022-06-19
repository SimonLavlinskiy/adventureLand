package actionsCounterController

import (
	"fmt"
	"project0/config"
	"project0/src/models"
)

func UserDo(user models.User, action string) {
	userAction := models.UserActionsCounter{
		UserId:     user.ID,
		ActionName: action,
		Count:      0,
	}
	userAction = GetOrCreateUserAction(userAction)

	userAction.Count += 1

	userAction.UpdateUserAction()
}

func GetOrCreateUserAction(action models.UserActionsCounter) models.UserActionsCounter {
	err := config.Db.
		Where(action).
		FirstOrCreate(&action).Error

	if err != nil {
		fmt.Println("GetOrCreateUserAction", err)
	}

	return action
}
