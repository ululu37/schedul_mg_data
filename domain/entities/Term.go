package entities

type Term struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"size:20" json:"name"`
}
