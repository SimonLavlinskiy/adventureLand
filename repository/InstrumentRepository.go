package repository

type Instrument struct {
	ID              uint `gorm:"primaryKey"`
	GoodId          *int `gorm:"embedded"`
	Good            *Item
	Type            string `gorm:"embedded"`
	CountDoing      int    `gorm:"embedded"`
	ItemsResultId   *int   `gorm:"embedded"`
	ItemsResult     *Item
	CountResultItem *int `gorm:"embedded"`
	NextStageItemId *int `gorm:"embedded"`
	NextStageItem   *Item
	Items           []Item `gorm:"many2many:instrument_item;"`
}
