package repository

type QuestTask struct {
	ID   uint   `gorm:"primaryKey"`
	Type string `gorm:"embedded"`
}
