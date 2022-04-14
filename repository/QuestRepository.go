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
	Result      Result
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
		Preload("Task").
		Preload("Result").
		Where(q).
		First(&results).Error

	if err != nil {
		fmt.Println("Нет такого квеста!")
	}

	return results
}

func (q Quest) QuestInfo(uq UserQuest) string {
	result := fmt.Sprintf("📜 *Задание* 📜\n`%s`\n\n"+
		"*Описание*: `%s`",
		q.Name, q.Description)

	if uq.Status != "" {
		result = fmt.Sprintf("%s\n\n*Статус*: _%s_",
			result, v.Get(fmt.Sprintf("quest_statuses.%s", uq.Status)))
	}

	return result
}
