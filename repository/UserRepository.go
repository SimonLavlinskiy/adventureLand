package repository

import (
	"database/sql"
	"fmt"
	"time"
)

type User struct {
	ID         uint `gorm:"primaryKey"`
	TgId       int
	TgChatId   int
	Username   string
	LocationId sql.NullString
	CreatedAt  time.Time
	Deleted    bool
}

func GetOrCreateUser(tgUserId int64) {
	fmt.Println(tgUserId)
}
