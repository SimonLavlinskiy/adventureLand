package models

import "time"

type UserSleep struct {
	UserId  uint
	SleptAt time.Time
}
