package entities

type Subject struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"size:200" json:"name"`
}
