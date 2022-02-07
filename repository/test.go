package repository

import uuid "github.com/satori/go.uuid"

type Test struct {
	ID          uuid.UUID `gorm:"primary_key;type:char(36);"`
	Name        string    `gorm:"embedded"`
	Description string    `gorm:"embedded"`
	SecTests    []Sectest `gorm:"ForeignKey:Uid"`
}
