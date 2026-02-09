package entities

type Curriculum struct {
	ID              uint   `gorm:"primaryKey" json:"id"`
	Name            string `gorm:"size:200" json:"name"`
	PreCurriculumID uint   `json:"pre_curriculum_id"`

	PreCurriculum PreCurriculum `gorm:"foreignKey:PreCurriculumID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"-"`

	SubjectInCurriculum []SubjectInCurriculum `gorm:"foreignKey:CurriculumID;" json:"subject_in_curriculums"`
}

type SubjectInCurriculum struct {
	ID                       uint  `gorm:"primaryKey" json:"id"`
	CurriculumID             uint  `json:"curriculum_id"`
	SubjectInPreCurriculumID uint  `json:"subject_in_pre_curriculum_id"`
	TermID                   *uint `json:"term_id"`

	Term                   Term                   `gorm:"foreignKey:TermID" json:"-"`
	SubjectInPreCurriculum SubjectInPreCurriculum `gorm:"foreignKey:SubjectInPreCurriculumID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE;" json:"-"`
	Curriculum             Curriculum             `gorm:"constraint:OnDelete:CASCADE,OnUpdate:CASCADE;" json:"-"`
}
