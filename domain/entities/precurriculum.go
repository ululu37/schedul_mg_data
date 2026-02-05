package entities

type PreCurriculum struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"size:200"`

	SubjectInPreCurriculum []SubjectInPreCurriculum `gorm:"foreignKey:PreCurriculumID;"`
}

type SubjectInPreCurriculum struct {
	ID              uint `gorm:"primaryKey"`
	PreCurriculumID uint
	SubjectID       uint
	Credit          int
	Subject         Subject

	PreCurriculum PreCurriculum `gorm:"constraint:OnDelete:CASCADE,OnUpdate:CASCADE;" json:"-"`
}
