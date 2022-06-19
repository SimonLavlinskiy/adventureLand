package models

type QuestTask struct {
	ID        uint   `gorm:"primaryKey"`
	Type      string `gorm:"embedded"`
	CountItem *int   `gorm:"embedded"`
	ItemId    *int   `gorm:"embedded"`
	Item      *Item
	UserAxisX *int `gorm:"embedded"`
	UserAxisY *int `gorm:"embedded"`
	MapId     *int `gorm:"embedded"`
}

func (t QuestTask) HasUserDoneTask(user User) bool {
	switch t.Type {
	case "haveItem":
		ui := UserItem{ItemId: *t.ItemId, UserId: int(user.ID)}.GetOrCreateUserItem()
		if *ui.Count >= *t.CountItem {
			return true
		}
	case "userLocation":
		uLoc := user.GetUserLocation()
		if uLoc.MapsId == t.MapId && uLoc.AxisX == t.UserAxisX && uLoc.AxisY == t.UserAxisY {
			return true
		}
	}

	return false
}
