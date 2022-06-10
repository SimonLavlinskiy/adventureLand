package actionsCounterController

import "project0/src/models"

func UserDo(user models.User, action string) {
	userAction := models.UserActionsCounter{
		UserId:     user.ID,
		ActionName: action,
		Count:      0,
	}
	userAction = models.GetOrCreateUserAction(userAction)

	userAction.Count += 1

	userAction.UpdateUserAction()
}
