package entities

type Term struct {
	ID   uint `gorm:"primaryKey"`
	Name string
}
