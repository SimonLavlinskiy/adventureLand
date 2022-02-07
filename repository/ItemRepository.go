package repository

type Item struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"embedded"`
	View string `gorm:"embedded"`
}
