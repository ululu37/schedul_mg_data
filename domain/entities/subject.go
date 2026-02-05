package entities

type Subject struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"size:200"`
}
