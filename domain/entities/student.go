package entities

type Student struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	AuthID       uint   `json:"auth_id"`
	Name         string `gorm:"size:200" json:"name"`
	CurriculumID *uint  `json:"curriculum_id"`
	Year         int    `json:"year"`
	ClassroomID  *uint  `json:"classroom_id"`

	Auth       Auth       `gorm:"constraint:OnDelete:CASCADE,OnUpdate:CASCADE;" json:"-"`
	Curriculum Curriculum `gorm:"foreignKey:CurriculumID;references:ID" json:"-"`
	Classroom  Classroom  `gorm:"foreignKey:ClassroomID;references:ID" json:"-"`
}
