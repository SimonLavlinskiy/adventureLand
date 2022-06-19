package models

import (
	"fmt"
	"project0/config"
	"time"
)

type UserQuest struct {
	ID        uint `gorm:"primaryKey"`
	UserId    uint `gorm:"embedded"`
	User      User
	QuestId   uint `gorm:"embedded"`
	Quest     Quest
	Status    string     `gorm:"embedded"`
	CreatedAt time.Time  `gorm:"autoCreateTime"`
	DoneAt    *time.Time `gorm:"embedded"`
}

func (uq UserQuest) GetUserQuest() UserQuest {
	var result UserQuest

	config.Db.
		Preload("User").
		Preload("Quest").
		Preload("Quest.Task").
		Preload("Quest.Result").
		Where(uq).
		First(&result)

	return result
}

func (uq UserQuest) CreateOrUpdateUserQuest() {
	result := UserQuest{
		UserId:  uq.UserId,
		QuestId: uq.QuestId,
		Status:  "processed",
	}

	err := config.Db.
		Preload("User").
		Preload("Quest").
		Preload("Quest.Task").
		Preload("Quest.Result").
		Where(uq).
		FirstOrCreate(&result).
		Error

	err = config.Db.Where(uq).Updates(UserQuest{Status: "processed"}).Error

	if err != nil {
		fmt.Println(err)
	}
}

func (uq UserQuest) UpdateUserQuest() bool {

	err := config.Db.
		Where(UserQuest{QuestId: uq.QuestId, UserId: uq.UserId}).
		Updates(UserQuest{
			Status: uq.Status,
			DoneAt: uq.DoneAt,
		}).
		Error

	if err != nil {
		return false
	}

	return true
}

func (uq UserQuest) UpdateUserQuestStatus() {
	err := config.Db.
		Where(UserQuest{
			QuestId: uq.QuestId,
			UserId:  uq.UserId,
		}).
		Updates(UserQuest{
			Status: uq.Status,
		}).
		Update("done_at", nil).
		Error

	if err != nil {
		fmt.Println("Update user quest status error:", err)
	}
}
