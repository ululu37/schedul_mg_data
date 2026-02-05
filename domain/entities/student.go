package entities

type Student struct {
	ID           uint `gorm:"primaryKey"`
	AuthID       uint
	Name         string `gorm:"size:200"`
	CurriculumID uint
	Year         int
	ClassroomID  uint

	Auth       Auth       `gorm:"constraint:OnDelete:CASCADE,OnUpdate:CASCADE;"`
	Curriculum Curriculum `gorm:"foreignKey:CurriculumID;references:ID"`
	Classroom  Classroom
}
