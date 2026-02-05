package entities

type Auth struct {
	ID        uint   `gorm:"primaryKey"`
	Username  string `gorm:"size:200;uniqueIndex;not null" json:"username"`
	Password  string `gorm:"size:200;not null"`
	HumanType string `gorm:"size:1" json:"human_type"` // 't' for teacher, 's' for student, '' for other
	Role      int    `gorm:"not null" json:"role"`     // 0: normal user, 1: admin
}

const (
	HumanTypeTeacher = "t" // teacher
	HumanTypeStudent = "s" // student
	HumanTypeOther   = ""  // other / unspecified
)

// Helper methods for HumanType
func (a *Auth) IsTeacher() bool { return a != nil && a.HumanType == HumanTypeTeacher }
func (a *Auth) IsStudent() bool { return a != nil && a.HumanType == HumanTypeStudent }
func (a *Auth) IsOther() bool   { return a == nil || a.HumanType == HumanTypeOther }

func HumanTypeDisplay(ht string) string {
	switch ht {
	case HumanTypeTeacher:
		return "teacher"
	case HumanTypeStudent:
		return "student"
	default:
		return "other"
	}
}

func (a *Auth) HumanTypeDisplay() string { return HumanTypeDisplay(a.HumanType) }
