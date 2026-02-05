package entities

type ScadulTeacher struct {
	ID        uint `gorm:"primaryKey"`
	TeacherID uint
	UseIn     string `gorm:"size:20"`

	Teacher                Teacher                  `gorm:"foreignKey:TeacherID"`
	SubjectInScadulTeacher []SubjectInScadulTeacher `gorm:"ForeignKey:ScadulTeacherID;"`
}

type SubjectInScadulTeacher struct {
	ID              uint `gorm:"primaryKey"`
	ScadulTeacherID uint
	TeacherID       uint
	SubjectID       uint
	Order           int

	ScadulTeacher ScadulTeacher `gorm:"constraint:OnDelete:CASCADE,OnUpdate:CASCADE;" json:"-"`
	Teacher       Teacher
	Subject       Subject
}
