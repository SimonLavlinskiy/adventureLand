package repository

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
