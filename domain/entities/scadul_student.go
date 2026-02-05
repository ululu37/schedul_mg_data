package entities

type ScadulStudent struct {
	ID                     uint `gorm:"primaryKey"`
	ClassroomID            uint
	UseIn                  string `gorm:"size:20"` // YYYY/term
	Classroom              Classroom
	SubjectInScadulStudent []SubjectInScadulStudent `gorm:"ForeignKey:ScadulStudentID;"`
}

type SubjectInScadulStudent struct {
	ID              uint `gorm:"primaryKey"`
	ScadulStudentID uint
	TeacherID       uint
	SubjectID       uint
	Order           int

	ScadulStudent ScadulStudent `gorm:"constraint:OnDelete:CASCADE,OnUpdate:CASCADE;" json:"-"`
	Teacher       Teacher
	Subject       Subject
}
