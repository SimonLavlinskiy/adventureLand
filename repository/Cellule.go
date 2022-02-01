package repository

type Cellule struct {
	ID      uint   `gorm:"primaryKey"`
	Map     string `gorm:"embedded"`
	AxisX   int    `gorm:"embedded"`
	AxisY   int    `gorm:"embedded"`
	View    string `gorm:"embedded"`
	CanStep bool   `gorm:"embedded"`
}
