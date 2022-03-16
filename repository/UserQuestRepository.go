package repository

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

func (uq UserQuest) GetOrCreateUserQuest() UserQuest {
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

	if err != nil {
		fmt.Println(err)
	}

	return result
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

func (uq UserQuest) UserDoneQuest(user User) bool {
	task := uq.Quest.Task

	switch task.Type {
	case "haveItem":
		ui := UserItem{ItemId: *task.ItemId, UserId: int(user.ID)}.UserGetUserItem()
		countItemResult := *ui.Count - *task.CountItem
		user.UpdateUserItem(UserItem{ID: ui.ID, Count: &countItemResult})
	case "userLocation":
		uLoc := user.GetUserLocation()
		if uLoc.MapsId == task.MapId && uLoc.AxisX == task.UserAxisX && uLoc.AxisY == task.UserAxisY {
			return true
		}
	}

	t := time.Now()
	uq.Status = "completed"
	uq.DoneAt = &t
	uq.UpdateUserQuest()

	return false
}
