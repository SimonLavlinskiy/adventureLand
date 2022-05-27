package models

type Teleport struct {
	ID     int `gorm:"primaryKey"`
	MapId  int `gorm:"embedded"`
	StartX int `gorm:"embedded"`
	StartY int `gorm:"embedded"`
}
