package entities

type PreCurriculum struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"size:200" json:"name"`

	SubjectInPreCurriculum []SubjectInPreCurriculum `gorm:"foreignKey:PreCurriculumID;" json:"subject_in_pre_curriculums"`
}

type SubjectInPreCurriculum struct {
	ID              uint    `gorm:"primaryKey" json:"id"`
	PreCurriculumID uint    `json:"pre_curriculum_id"`
	SubjectID       uint    `json:"subject_id"`
	Credit          int     `json:"credit"`
	Subject         Subject `gorm:"foreignKey:SubjectID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE;" json:"subject"`

	PreCurriculum PreCurriculum `gorm:"constraint:OnDelete:CASCADE,OnUpdate:CASCADE;" json:"-"`
}
