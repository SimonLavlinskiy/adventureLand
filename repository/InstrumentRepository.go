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

func GetInstrumentsUserCanUse(user User, cell Cellule) []string {
	var instrumentsUserCanUse []string
	instruments := cell.Item.Instruments

	for _, instrument := range instruments {
		if user.LeftHandId != nil && user.LeftHand.Type == instrument.Good.Type {
			instrumentsUserCanUse = append(instrumentsUserCanUse, user.LeftHand.View)
		}
		if user.RightHandId != nil && user.RightHand.Type == instrument.Good.Type {
			instrumentsUserCanUse = append(instrumentsUserCanUse, user.RightHand.View)
		}
		if user.HeadId != nil && user.Head.Type == instrument.Good.Type {
			instrumentsUserCanUse = append(instrumentsUserCanUse, user.Head.View)
		}
	}
	if cell.Item.CanTake {
		instrumentsUserCanUse = append(instrumentsUserCanUse, "👋")
	}
	if cell.CanStep {
		instrumentsUserCanUse = append(instrumentsUserCanUse, "\U0001F9B6")
	}

	return instrumentsUserCanUse
}
