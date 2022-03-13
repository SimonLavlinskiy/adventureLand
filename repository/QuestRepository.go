package repository

import (
	"fmt"
	v "github.com/spf13/viper"
	"project0/config"
)

type Quest struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"embedded"`
	Description string `gorm:"embedded"`
	Type        string `gorm:"embedded"`
	ResultId    int    `gorm:"embedded"`
	Result      QuestResult
	TaskId      int `gorm:"embedded"`
	Task        QuestTask
}

func (q Quest) GetQuests() []Quest {
	var results []Quest

	err := config.Db.
		Find(&results).Error

	if err != nil {
		fmt.Println("Нет квестов!")
	}

	return results
}

func (q Quest) GetQuest() Quest {
	var results Quest

	err := config.Db.
		Where(q).
		First(&results).Error

	if err != nil {
		fmt.Println("Нет такого квеста!")
	}

	return results
}

func (q Quest) GetUserQuest() *UserQuest {
	var result UserQuest

	err := config.Db.
		Preload("Quest").
		Preload("Quest.Task").
		Preload("Quest.Result").
		Where(UserQuest{QuestId: q.ID}).
		Find(&result).
		Error

	if err != nil {
		return nil
	}

	return &result
}

func (q Quest) QuestInfo(uq *UserQuest) string {
	result := fmt.Sprintf("*Название*: _%s_\n"+
		"*Описание*: _%s_",
		q.Name, q.Description)

	if uq != nil {
		result = fmt.Sprintf("_%s_\n"+
			"*Статус:*: _%s_",
			result, v.Get(fmt.Sprintf("quest_statuses.%s", uq.Status)))
	}

	return result
}
