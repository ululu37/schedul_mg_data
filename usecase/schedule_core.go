package usecase

import (
	"errors"
	"scadulDataMono/domain/entities"
)

type SchedulRes struct {
	studentScheduls []entities.ScadulStudent
	teacherScheduls []entities.ScadulTeacher
}

type TeacherAssign struct {
	TeacherID     uint
	AllocatedHour int
	subjects      []SubjectAssign
}
type SubjectAssign struct {
	SubjectInPreCurriculumID uint
	Credit                   int
	ClassroomID              uint
}

func Can(teachers []entities.Teacher, classrooms []entities.Classroom) (SchedulRes, error) {
	//	allSubject := GetAllSubjectInClassroom(classrooms)
	//	teacherAssign := ConvertTeacherAssign(teachers)
	everageHour := EverageHour(teachers, classrooms) //หาค่าเฉลีย
	if everageHour == 0 || everageHour > 40 {
		return SchedulRes{}, errors.New("ไม่สามารถจัดตารางได้")
	}

	return SchedulRes{}, nil
}

func GetAllSubjectInClassroom(classrooms []entities.Classroom) []entities.SubjectInPreCurriculum {
	allSubject := []entities.SubjectInPreCurriculum{}
	for _, cl := range classrooms {
		for _, subject := range cl.PreCurriculum.SubjectInPreCurriculum {
			allSubject = append(allSubject, subject)
		}
	}
	return allSubject
}

func ConvertTeacherAssign(teachers []entities.Teacher) []TeacherAssign {
	teacherAssign := []TeacherAssign{}
	for _, teacher := range teachers {
		teacherAssign = append(teacherAssign, TeacherAssign{
			TeacherID:     teacher.ID,
			AllocatedHour: 0,
			subjects:      []SubjectAssign{},
		})
	}
	return teacherAssign
}

func EverageHour(teachers []entities.Teacher, classrooms []entities.Classroom) int { //หาค่าเฉลีย
	allSubject := []entities.SubjectInPreCurriculum{}
	allHour := 0

	for _, cl := range classrooms {
		for _, subject := range cl.PreCurriculum.SubjectInPreCurriculum {
			allSubject = append(allSubject, subject)
			allHour += subject.Credit
		}
	}

	return allHour / len(teachers)
}
