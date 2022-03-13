package repository

import "project0/config"

type UserQuest struct {
	UserId  uint `gorm:"embedded"`
	User    User
	QuestId uint `gorm:"embedded"`
	Quest   Quest
	Status  string `gorm:"embedded"`
}

func (uq UserQuest) GetUserQuest() *UserQuest {
	var result *UserQuest

	err := config.Db.
		Preload("Quest").
		Preload("Quest.Task").
		Preload("Quest.Result").
		Where(uq).
		Find(result).
		Error

	if err != nil {
		return nil
	}

	return result
}
