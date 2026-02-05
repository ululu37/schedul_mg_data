package entities

type Teacher struct {
	ID     uint   `gorm:"primaryKey"`
	Name   string `gorm:"size:200" json:"name"`
	Resume string `gorm:"type:text" json:"resume"`
	AuthID uint   `json:"auth_id"`

	Auth      Auth               `gorm:"constraint:OnDelete:CASCADE,OnUpdate:CASCADE; " json:"-"`
	MySubject []TeacherMySubject `gorm:"foreignKey:TeacherID" json:"-"`
}

type TeacherMySubject struct {
	ID         uint `gorm:"primaryKey"`
	TeacherID  uint `่json:"teacher_id"`
	SubjectID  uint `json:"subject_id"`
	Preference int  `json:"preference"`

	Subject Subject
	Teacher Teacher `gorm:"constraint:OnDelete:CASCADE,OnUpdate:CASCADE;" json:"-"`
}
