package resultController

import (
	"fmt"
	"project0/src/controllers/userController"
	"project0/src/models"
	"project0/src/repositories"
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

	if r.Money != nil {
		userMoney := *u.Money + *r.Money
		repositories.UpdateUser(models.User{TgId: u.TgId, Money: &userMoney})
	}
}

func UserGetResultMsg(result models.Result) string {
	result = result.GetResult()

	msg := "🏆 *Ты получил*:"
	if result.Item != nil {
		msg = fmt.Sprintf("%s\n_%s %s - %d шт._", msg, result.Item.View, result.Item.Name, *result.CountItem)
	}
	if result.SpecialItem != nil {
		msg = fmt.Sprintf("%s\n_%s %s - %d шт._", msg, result.SpecialItem.View, result.SpecialItem.Name, *result.SpecialItemCount)
	}
	if result.Money != nil {
		msg = fmt.Sprintf("%s\n_💰 %d 💰_", msg, *result.Money)
	}

	return msg
}
