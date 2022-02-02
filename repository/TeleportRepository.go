package repository

import "project0/config"

type Teleport struct {
	ID     int    `gorm:"primaryKey"`
	Map    string `gorm:"embedded"`
	StartX int    `gorm:"embedded"`
	StartY int    `gorm:"embedded"`
}

func GetTeleport(teleport Teleport) Teleport {
	result := Teleport{}

	err := config.Db.Where(teleport).FirstOrCreate(&result).Error

	if err != nil {
		panic(err)
	}

	return result
}
