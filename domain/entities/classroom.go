package entities

type Classroom struct {
	ID              uint   `gorm:"primaryKey" json:"id"`
	Name            string `gorm:"size:200" json:"name"`
	PreCurriculumID *uint  `json:"pre_curriculum_id"`
	CurriculumID    *uint  `json:"curriculum_id"`
	Year            *int   `json:"year"`

	Student       []Student     `gorm:"ForeignKey:ClassroomID;references:ID;" json:"students"`
	PreCurriculum PreCurriculum `gorm:"constraint:OnDelete:CASCADE,OnUpdate:CASCADE;" json:"pre_curriculum"`
	Curriculum    Curriculum    `gorm:"constraint:OnDelete:CASCADE,OnUpdate:CASCADE;" json:"curriculum"`
}
