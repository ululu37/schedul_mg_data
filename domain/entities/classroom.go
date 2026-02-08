package entities

type Classroom struct {
	ID              uint          `gorm:"primaryKey"`
	Name            string        `gorm:"size:200"`
	PreCurriculumID *uint         `json:"pre_curriculum_id"`
	CurriculumID    *uint         `json:"curriculum_id"`
	Year            *int          `json:"year"`
	Student         []Student     `gorm:"ForeignKey:ClassroomID;"`
	PreCurriculum   PreCurriculum `gorm:"constraint:OnDelete:CASCADE,OnUpdate:CASCADE;" json:"-"`
	Curriculum      Curriculum    `gorm:"constraint:OnDelete:CASCADE,OnUpdate:CASCADE;" json:"-"`
}
