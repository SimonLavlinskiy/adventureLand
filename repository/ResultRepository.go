package repository

import "project0/config"

type Result struct {
	ID               uint   `gorm:"primaryKey"`
	Type             string `gorm:"embedded"`
	ItemId           *uint  `gorm:"embedded"`
	Item             *Item
	CountItem        *uint `gorm:"embedded"`
	SpecialItemId    *uint `gorm:"embedded"`
	SpecialItem      *Item
	SpecialItemCount *uint `gorm:"embedded"`
	Experience       *int  `gorm:"embedded"`
}

func (u User) UserGetResult(r Result) {
	switch r.Type {
	case "casual":
		u.UserGetExperience(r)
	case "casualPlus":
		u.UserGetExperience(r)
		u.UserGetResultItem(r)
	case "superCasual":
		u.UserGetExperience(r)
		u.UserGetResultItem(r)
		u.UserGetResultSpecialItem(r)
	}
}

func (r Result) GetResult() Result {
	var res Result
	config.Db.
		Preload("Item").
		Preload("SpecialItem").
		Where(Result{ID: r.ID}).
		First(&res)

	return res
}
