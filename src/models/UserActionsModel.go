package models

import (
	"fmt"
	"project0/config"
)

type UserActionsCounter struct {
	ID         uint   `gorm:"primaryKey"`
	UserId     uint   `gorm:"embedded"`
	ActionName string `gorm:"embedded"`
	Count      uint   `gorm:"embedded"`
}

func GetOrCreateUserAction(action UserActionsCounter) UserActionsCounter {
	err := config.Db.
		Where(action).
		FirstOrCreate(&action).Error

	if err != nil {
		fmt.Println("GetOrCreateUserAction", err)
	}

	return action
}

func (u UserActionsCounter) UpdateUserAction() {
	err := config.Db.
		Updates(&u).
		Error

	if err != nil {
		fmt.Println("UpdateUserAction error: ", err)
	}
}

func (u UserActionsCounter) GetStepsPlace() int {
	var actions []UserActionsCounter
	config.Db.
		Where("action_name = ?", "step").
		Where("count >= ?", u.Count).
		Find(&actions)

	return len(actions)
}
