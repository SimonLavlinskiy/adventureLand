package resultController

import (
	"project0/src/controllers/userController"
	"project0/src/models"
)

func UserGetResult(u models.User, r models.Result) {
	switch r.Type {
	case "casual":
		userController.UserGetExperience(u, r)
	case "casualPlus":
		userController.UserGetExperience(u, r)
		u.UserGetResultItem(r)
	case "superCasual":
		userController.UserGetExperience(u, r)
		u.UserGetResultItem(r)
		u.UserGetResultSpecialItem(r)
	}
}
