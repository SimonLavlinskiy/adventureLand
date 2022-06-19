package userController

import (
	"fmt"
	v "github.com/spf13/viper"
	"project0/src/controllers/actionsCounterController"
	"project0/src/models"
	"project0/src/repositories"
	"project0/src/services/helpers"
)

func UpdateUserHand(user models.User, char []string) (models.User, models.UserItem) {
	userItem := models.UserItem{ID: helpers.ToInt(char[1])}.GetOrCreateUserItem()

	switch char[0] {
	case v.GetString("callback_char.change_left_hand"):
		clothes := &models.Clothes{LeftHandId: &userItem.ItemId}
		repositories.UpdateUser(models.User{TgId: user.TgId, Clothes: *clothes})

	case v.GetString("callback_char.change_right_hand"):
		clothes := &models.Clothes{RightHandId: &userItem.ItemId}
		repositories.UpdateUser(models.User{TgId: user.TgId, Clothes: *clothes})
	}

	user = repositories.GetUser(models.User{TgId: user.TgId})

	return user, userItem
}

func UserGetExperience(user models.User, r models.Result) {
	resultExp := user.Experience + *r.Experience
	repositories.UpdateUser(models.User{ID: user.ID, Experience: resultExp})
}

func UserBuyHome(u models.User, m models.Map) {
	*u.Money -= v.GetInt("main_info.cost_of_house")
	u.HomeId = &m.ID

	repositories.UpdateUser(u)
}

func GetUserInfo(u models.User) string {
	step := actionsCounterController.GetOrCreateUserAction(models.UserActionsCounter{UserId: u.ID, ActionName: "step"})
	stepPlace := step.GetStepsPlace()
	messageMap := fmt.Sprintf("🔅 🔆 *Профиль* 🔆 🔅\n\n"+
		"*Твое имя*: %s\n"+
		"*Аватар*: %s\n"+
		"*Золото*: %d 💰\n"+
		"*Здоровье*: _%d_ ❤️\n"+
		"*Сытость*: _%d_ 😋️\n"+
		"*Шаги*: _%d_ 👣 (_%d место_)",
		u.Username, u.Avatar, *u.Money, u.Health, u.Satiety, step.Count, stepPlace)

	return messageMap
}
