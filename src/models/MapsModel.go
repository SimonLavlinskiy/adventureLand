package models

import (
	"fmt"
	"project0/config"
)

type Map struct {
	ID               uint   `gorm:"primaryKey"`
	Name             string `gorm:"embedded"`
	SizeX            int    `gorm:"embedded"`
	SizeY            int    `gorm:"embedded"`
	StartX           int    `gorm:"embedded"`
	StartY           int    `gorm:"embedded"`
	DayType          string `gorm:"embedded"`
	EmptySpaceSymbol string `gorm:"embedded"`
}

func (m Map) CreateMap() Map {
	err := config.Db.
		Create(&m).
		Error

	if err != nil {
		fmt.Println(err)
	}

	return m
}

type UserMap struct {
	LeftIndent  int
	RightIndent int
	UpperIndent int
	DownIndent  int
}

type MapButtons struct {
	Up         string
	UpData     string
	Left       string
	LeftData   string
	Right      string
	RightData  string
	Down       string
	DownData   string
	Center     string
	CenterData string
}
