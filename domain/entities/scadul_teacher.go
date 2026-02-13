package entities

type ScadulTeacher struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	TeacherID uint   `json:"teacher_id"`
	UseIn     string `gorm:"size:20" json:"use_in"`

	Teacher                Teacher                  `gorm:"foreignKey:TeacherID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE;" json:"teacher"`
	SubjectInScadulTeacher []SubjectInScadulTeacher `gorm:"ForeignKey:ScadulTeacherID;" json:"subject_in_scadul_teachers"`
}

type SubjectInScadulTeacher struct {
	ID              uint `gorm:"primaryKey" json:"id"`
	ScadulTeacherID uint `json:"scadul_teacher_id"`
	ClassroomID     uint `json:"classroom_id"`
	SubjectID       uint `json:"subject_id"`
	Order           int  `json:"order"`

	ScadulTeacher ScadulTeacher `gorm:"foreignKey:ScadulTeacherID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE;" json:"-"`
	Classroom     Classroom     `gorm:"foreignKey:ClassroomID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE;" json:"classroom"`
	Subject       Subject       `gorm:"foreignKey:SubjectID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE;" json:"subject"`
}
