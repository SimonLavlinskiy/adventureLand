package repository

type QuestResult struct {
	ID   uint   `gorm:"primaryKey"`
	Type string `gorm:"embedded"`
}
