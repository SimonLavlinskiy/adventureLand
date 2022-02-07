package repository

import (
	_ "fmt"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	uuid "github.com/satori/go.uuid"
)

type Sectest struct {
	ID             uuid.UUID `gorm:"primary_key;type:char(36);"`
	NameSec        string
	DescriptionSec string
	Uid            uuid.UUID
}
