package entities

type Classroom struct {
	ID        uint `gorm:"primaryKey"`
	TeacherID uint
	Year      int
	Teacher Teacher
	Student []Student `gorm:"ForeignKey:ClassroomID;"`
}
