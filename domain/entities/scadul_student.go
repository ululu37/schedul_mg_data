package entities

type ScadulStudent struct {
	ID                     uint                     `gorm:"primaryKey" json:"id"`
	ClassroomID            uint                     `json:"classroom_id"`
	UseIn                  string                   `gorm:"size:20" json:"use_in"` // YYYY/term
	Classroom              Classroom                `gorm:"foreignKey:ClassroomID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE;" json:"-"`
	SubjectInScadulStudent []SubjectInScadulStudent `gorm:"ForeignKey:ScadulStudentID;" json:"subject_in_scadul_students"`
}

type SubjectInScadulStudent struct {
	ID              uint `gorm:"primaryKey" json:"id"`
	ScadulStudentID uint `json:"scadul_student_id"`
	TeacherID       uint `json:"teacher_id"`
	SubjectID       uint `json:"subject_id"`
	Order           int  `json:"order"`

	ScadulStudent ScadulStudent `gorm:"constraint:OnDelete:CASCADE,OnUpdate:CASCADE;" json:"-"`
	Teacher       Teacher       `gorm:"foreignKey:TeacherID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE;" json:"-"`
	Subject       Subject       `gorm:"foreignKey:SubjectID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE;" json:"-"`
}
