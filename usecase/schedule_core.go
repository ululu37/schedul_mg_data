package usecase

import (
	"errors"
	"fmt"
	"scadulDataMono/domain/entities"
)

type SchedulRes struct {
	StudentScheduls []entities.ScadulStudent
	TeacherScheduls []entities.ScadulTeacher
}

type TeacherAssign struct {
	AllocatedHour int
	Subjects      map[string]SubjectAssign           //[SubjectID-ClassroomID]SubjectAssign
	MySubjects    map[uint]entities.TeacherMySubject //[SubjectID]
}
type SubjectAssign struct {
	SubjectID   uint
	Credit      int
	ClassroomID uint
}
type position struct {
	x int
	y int
}
type ResultAssign struct {
	teacherAssign   map[uint]TeacherAssign
	classroomAssign map[string]struct{}
	erverageHour    int
}

func AssignSubject(everageHour int, teacherAssign map[uint]TeacherAssign, classroomAssign map[string]struct{}, teachers []entities.Teacher, classrooms []entities.Classroom, is_everage bool) (*ResultAssign, error) {

	if everageHour == 0 || everageHour > 40 {
		return nil, errors.New("ไม่สามารถจัดตารางได้")
	}

	// เริ่มเเจกงาน
	for _, cl := range classrooms {
		for _, subjectInCurriculum := range cl.Curriculum.SubjectInCurriculum {
			subject := subjectInCurriculum.SubjectInPreCurriculum
			if !is_everage {
				if _, ok := classroomAssign[fmt.Sprintf("%d-%d", cl.ID, subject.SubjectID)]; ok {
					fmt.Println("skip", cl.ID, subject.SubjectID)
					continue
				}
			}

			//หาครูที่สอนวิชา นี้ ได้
			teacherTech := make(map[uint]int) // [TeacherID]preferance
			for id, teacher := range teacherAssign {
				if _, ok := teacher.MySubjects[subject.SubjectID]; ok {

					if teacher.AllocatedHour+subject.Credit > 40 {

						continue
					}
					if teacher.AllocatedHour+subject.Credit > everageHour && is_everage {

						continue
					}

					teacherTech[id] = teacher.MySubjects[subject.SubjectID].Preference

				}
			}
			isBreak := false
			//เอาคนคะเเนน เยอะก่อน
			for i := 10; i != 1; i-- {
				///fmt.Println(i)
				//คัดวิชา prefreranc i
				mostPreferanceteacher := []uint{}
				for id, preferance := range teacherTech {
					if preferance == i {

						mostPreferanceteacher = append(mostPreferanceteacher, id)
					}
				}
				//fmt.Println(mostPreferanceteacher)
				//เรียงตามจำนวนชั่วโมง น้อยไปมาก
				if len(mostPreferanceteacher) == 0 {
					continue
				}

				if len(mostPreferanceteacher) > 1 {
					for r := 1; r < len(mostPreferanceteacher); r++ {

						for it := range mostPreferanceteacher {
							if it == len(mostPreferanceteacher)-r {

								break
							}
							if teacherAssign[mostPreferanceteacher[it]].AllocatedHour > teacherAssign[mostPreferanceteacher[it+1]].AllocatedHour {
								cache := mostPreferanceteacher[it]
								mostPreferanceteacher[it] = mostPreferanceteacher[it+1]
								mostPreferanceteacher[it+1] = cache
							}

						}
					}
				}

				//เเจกงาน
				for _, teacher := range mostPreferanceteacher {
					if teacherAssign[teacher].AllocatedHour+subject.Credit > 40 {
						continue
					}

					// Get the struct from the map
					ta := teacherAssign[teacher]
					// Modify the struct

					ta.AllocatedHour += subject.Credit

					// Initialize the Subjects map if it's nil
					//fmt.Println(cl.CurriculumID)
					key := fmt.Sprintf("%d-%d", subject.SubjectID, cl.ID)
					ta.Subjects[key] = SubjectAssign{
						SubjectID:   subject.SubjectID,
						Credit:      subject.Credit,
						ClassroomID: cl.ID,
					}
					// Put the modified struct back into the map
					teacherAssign[teacher] = ta
					classroomAssign[fmt.Sprintf("%d-%d", cl.ID, subject.SubjectID)] = struct{}{}
					isBreak = true

					break

				}
				if isBreak {
					break
				}
			}
			if !is_everage {
				if !isBreak {
					return nil, errors.New("ไม่สามารถจัดตารางได้")
				}
			}
		}

	}
	return &ResultAssign{
		teacherAssign:   teacherAssign,
		classroomAssign: classroomAssign,
		erverageHour:    everageHour,
	}, nil

}

func Can(teachers []entities.Teacher, classrooms []entities.Classroom) (SchedulRes, error) {
	fmt.Println("classrooms", len(classrooms))
	for _, cl := range classrooms {
		fmt.Println("cl", cl.Name)
		for i, subject := range cl.Curriculum.SubjectInCurriculum {
			fmt.Printf("%d: %s\n", i, subject.SubjectInPreCurriculum.Subject.Name)
		}
	}
	c := make(map[string]struct{})
	a := ConvertTeacherAssign(teachers)
	everageHour := EverageHour(teachers, classrooms) //หาค่าเฉลีย
	fmt.Println("everageHour", everageHour)

	var result *ResultAssign
	res, err := AssignSubject(everageHour, a, c, teachers, classrooms, true)
	result = res
	if err != nil {
		return SchedulRes{}, err
	}
	//fmt.Println(result.classroomAssign)
	// for id, teacher := range result.teacherAssign {
	// 	fmt.Printf("teacher %d : %d\n", id, teacher.AllocatedHour)
	// }
	//fmt.Println("รอบสอง")
	res2, err := AssignSubject(everageHour, result.teacherAssign, result.classroomAssign, teachers, classrooms, false)
	result = res2
	if err != nil {
		return SchedulRes{}, err
	}

	for id, teacher := range result.teacherAssign {
		fmt.Printf("teacher %d : %d\n", id, teacher.AllocatedHour)
	}
	oderTable := orderSchedule(result.teacherAssign, classrooms)
	for i := 1; i <= 40; i++ {
		fmt.Printf("%d : %v %v %v\n", i, oderTable[position{x: 0, y: i}], oderTable[position{x: 1, y: i}], oderTable[position{x: 2, y: i}])
	}

	// Build SchedulRes
	studentScheduleMap := make(map[uint]*entities.ScadulStudent)
	teacherScheduleMap := make(map[uint]*entities.ScadulTeacher)

	// Initialize maps
	for _, cl := range classrooms {
		studentScheduleMap[cl.ID] = &entities.ScadulStudent{
			ClassroomID:            cl.ID,
			SubjectInScadulStudent: []entities.SubjectInScadulStudent{},
			UseIn:                  "2024/1", // Example default
		}
	}
	for _, t := range teachers {
		teacherScheduleMap[t.ID] = &entities.ScadulTeacher{
			TeacherID:              t.ID,
			SubjectInScadulTeacher: []entities.SubjectInScadulTeacher{},
			UseIn:                  "2024/1", // Example default
		}
	}

	for pos, assign := range oderTable {
		if assign == nil {
			continue
		}
		// Classroom schedule
		clID := classrooms[pos.x].ID
		if ss, ok := studentScheduleMap[clID]; ok {
			ss.SubjectInScadulStudent = append(ss.SubjectInScadulStudent, entities.SubjectInScadulStudent{
				TeacherID: assign.TeacherID,
				SubjectID: assign.SubjectID,
				Order:     pos.y,
			})
		}

		// Teacher schedule
		if ts, ok := teacherScheduleMap[assign.TeacherID]; ok {
			ts.SubjectInScadulTeacher = append(ts.SubjectInScadulTeacher, entities.SubjectInScadulTeacher{
				ClassroomID: assign.ClassroomID,
				SubjectID:   assign.SubjectID,
				Order:       pos.y,
			})
		}
	}

	studentRes := []entities.ScadulStudent{}
	for _, v := range studentScheduleMap {
		studentRes = append(studentRes, *v)
	}
	teacherRes := []entities.ScadulTeacher{}
	for _, v := range teacherScheduleMap {
		teacherRes = append(teacherRes, *v)
	}

	return SchedulRes{
		StudentScheduls: studentRes,
		TeacherScheduls: teacherRes,
	}, nil
}

func orderSchedule(teacherAssign map[uint]TeacherAssign, classrooms []entities.Classroom) map[position]*TeacherInOrder {
	teacher := NewTeacherAllocate(teacherAssign)
	order := make(map[position]*TeacherInOrder)
	for i := 1; i <= 40; i++ {
		for roomI := range classrooms {
			//fmt.Println(i, roomI)
			order[position{x: roomI, y: i}] = teacher.tik(i, classrooms[roomI].ID)
		}
	}
	//fmt.Println("order", order)
	return order
}

type TeacherInOrder struct {
	TeacherID   uint
	SubjectID   uint
	ClassroomID uint
}

type teacherAllocate struct {
	TeacherAssign map[uint]TeacherAssign
	YAdded        map[int]InAdded
}
type InAdded struct {
	teachers map[uint]struct{}
}

func NewTeacherAllocate(TeacherAssign map[uint]TeacherAssign) teacherAllocate {
	yAdded := make(map[int]InAdded)
	for i := 1; i <= 40; i++ {
		teachers := make(map[uint]struct{})
		for teacherID := range TeacherAssign {
			teachers[teacherID] = struct{}{}
		}
		yAdded[i] = InAdded{teachers: teachers}
	}
	return teacherAllocate{
		TeacherAssign: TeacherAssign,
		YAdded:        yAdded,
	}
}

func (ta *teacherAllocate) tik(y int, classroomID uint) *TeacherInOrder {
	var winID uint
	var winSubjectID uint
	var winKey string
	var maxHour int = -1
	found := false

	// หาครูที่มีชั่วโมงมากที่สุดและยังไม่ได้ถูกเพิ่มในคาบ y
	if len(ta.YAdded[y].teachers) > 0 {

		for teacherID := range ta.YAdded[y].teachers {
			// Check if teacher has subjects for this classroom
			t := ta.TeacherAssign[teacherID]
			for key, sub := range t.Subjects {
				if sub.ClassroomID == classroomID && sub.Credit > 0 {
					if t.AllocatedHour > maxHour {
						maxHour = t.AllocatedHour
						winID = teacherID
						winSubjectID = sub.SubjectID
						winKey = key
						found = true
					}
					// Found a valid subject for this teacher, compare with maxHour
					// We break inner loop because any valid subject qualifies the teacher
					// But wait, if teacher has multiple subjects for same class?
					// Just picking one is fine.
					break
				}
			}
		}
	} else {
		return nil
	}
	if !found {
		return nil
	}

	delete(ta.YAdded[y].teachers, winID)
	//fmt.Printf("winID %v\n", ta.YAdded[y].teachers)
	// Get the struct from the map, modify it, and put it back
	teacher := ta.TeacherAssign[winID]
	teacher.AllocatedHour -= 1

	// Decrement subject credit
	if s, ok := teacher.Subjects[winKey]; ok {
		s.Credit -= 1
		teacher.Subjects[winKey] = s
	}

	ta.TeacherAssign[winID] = teacher
	//fmt.Printf("winID %d\n", winID)
	//fmt.Printf("%d %d\n", ta.TeacherAssign[1].AllocatedHour, ta.TeacherAssign[2].AllocatedHour)

	return &TeacherInOrder{TeacherID: winID, SubjectID: winSubjectID, ClassroomID: classroomID}
}

func GetAllSubjectInClassroom(classrooms []entities.Classroom) []entities.SubjectInPreCurriculum {

	allSubject := []entities.SubjectInPreCurriculum{}
	for _, cl := range classrooms {
		for _, subject := range cl.Curriculum.SubjectInCurriculum {
			allSubject = append(allSubject, subject.SubjectInPreCurriculum)
		}
	}
	return allSubject
}

func ConvertTeacherAssign(teachers []entities.Teacher) map[uint]TeacherAssign {
	teacherAssign := make(map[uint]TeacherAssign, len(teachers))
	for _, teacher := range teachers {
		teacherAssign[teacher.ID] = TeacherAssign{

			AllocatedHour: 0,
			Subjects:      map[string]SubjectAssign{},
			MySubjects:    ConvertMySubject(teacher.MySubject),
		}
	}
	return teacherAssign
}

func ConvertMySubject(mySubjects []entities.TeacherMySubject) map[uint]entities.TeacherMySubject {
	mySubject := make(map[uint]entities.TeacherMySubject, len(mySubjects))
	for _, ms := range mySubjects {
		mySubject[ms.SubjectID] = ms
	}
	return mySubject
}

func EverageHour(teachers []entities.Teacher, classrooms []entities.Classroom) int { //หาค่าเฉลีย
	allSubject := []entities.SubjectInPreCurriculum{}
	allHour := 0

	for _, cl := range classrooms {
		for _, subject := range cl.Curriculum.SubjectInCurriculum {
			allSubject = append(allSubject, subject.SubjectInPreCurriculum)
			allHour += subject.SubjectInPreCurriculum.Credit
		}
	}

	return allHour / len(teachers)
}
