package models

type Receipt struct {
	ID              uint `gorm:"primaryKey"`
	Component1ID    *int `gorm:"embedded"`
	Component1      *Item
	Component1Count *int `gorm:"embedded"`
	Component2ID    *int `gorm:"embedded"`
	Component2      *Item
	Component2Count *int `gorm:"embedded"`
	Component3ID    *int `gorm:"embedded"`
	Component3      *Item
	Component3Count *int `gorm:"embedded"`
	ItemResultID    int
	ItemResult      Item
	ItemResultCount *int `gorm:"embedded"`
}
