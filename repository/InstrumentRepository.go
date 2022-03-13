package repository

type Instrument struct {
	ID                 uint `gorm:"primaryKey"`
	GoodId             *int `gorm:"embedded"`
	Good               *Item
	Type               string `gorm:"embedded"`
	ItemsResultId      *int   `gorm:"embedded"`
	ItemsResult        *Item
	CountResultItem    *int `gorm:"embedded"`
	NextStageItemId    *int `gorm:"embedded"`
	NextStageItem      *Item
	CountNextStageItem *int   `gorm:"embedded"`
	Items              []Item `gorm:"many2many:instrument_item;"`
}

type Y struct {
	i int
	s string
}

func GetInstrumentsUserCanUse(user User, cell Cell) map[string]string {
	instrumentsUserCanUse := map[string]string{}
	instruments := cell.Item.Instruments

	for _, instrument := range instruments {
		if user.LeftHandId != nil && user.LeftHand.ID == instrument.Good.ID {
			instrumentsUserCanUse[user.LeftHand.View] = user.LeftHand.Type
		}
		if user.RightHandId != nil && user.RightHand.ID == instrument.Good.ID {
			instrumentsUserCanUse[user.RightHand.View] = user.RightHand.Type
		}
		if user.HeadId != nil && user.Head.ID == instrument.Good.ID {
			instrumentsUserCanUse[user.Head.View] = user.Head.Type
		}
	}
	if cell.Item.CanTake {
		instrumentsUserCanUse["👋"] = "hand"
	}
	if cell.Item.CanStep && *cell.Type != "swap" {
		instrumentsUserCanUse["\U0001F9B6"] = "step"
	}

	return instrumentsUserCanUse
}
