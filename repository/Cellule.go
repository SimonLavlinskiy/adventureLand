package repository

type Cellule struct {
	ID      uint   `gorm:"primaryKey"`
	Map     string `gorm:"embedded"`
	AxisX   uint   `gorm:"embedded"`
	AxisY   uint   `gorm:"embedded"`
	View    string `gorm:"embedded"`
	CanStep bool   `gorm:"embedded"`
}
