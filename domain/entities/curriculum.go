package entities

type Curriculum struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"size:200"`

	SubjectInCurriculum []SubjectInCurriculum `gorm:"foreignKey:CurriculumID;"`
}

type SubjectInCurriculum struct {
	ID                       uint `gorm:"primaryKey"`
	CurriculumID             uint
	SubjectInPreCurriculumID uint
	TermID                   uint

	Term                   Term                   `gorm:"foreignKey:TermID"`
	SubjectInPreCurriculum SubjectInPreCurriculum `gorm:"foreignKey:SubjectInPreCurriculumID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE;"`
	Curriculum             Curriculum             `gorm:"constraint:OnDelete:CASCADE,OnUpdate:CASCADE;" json:"-"`
}

type TermInCurriculum struct {
	ID     uint `gorm:"primaryKey"`
	TermID uint
}
