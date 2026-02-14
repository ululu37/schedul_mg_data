package usecase

import (
	"errors"
	"fmt"
	"scadulDataMono/domain/entities"
	"sort"
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

type TeacherRule struct {
	TeacherID        uint
	Consecutive      int   // เช่น 4 คาบติด
	BusyPeriods      []int // คาบที่ไม่ว่างจริงๆ
	PreferredPeriods []int // คาบที่อยากสอนเป็นพิเศษ
}

func AssignSubject(everageHour int, teacherAssign map[uint]TeacherAssign, classroomAssign map[string]struct{}, teachers []entities.Teacher, classrooms []entities.Classroom, is_everage bool) (*ResultAssign, error) {

	if everageHour == 0 || everageHour > 40 {
		return nil, errors.New("ไม่สามารถจัดตารางได้")
	}

	// 1. ดึงวิชาจากทุกห้องออกมารวมกัน (Global Pool)
	type GlobalSubject struct {
		ClassroomID uint
		Subject     entities.SubjectInPreCurriculum
	}
	var allSubjects []GlobalSubject
	for _, cl := range classrooms {
		for _, sic := range cl.Curriculum.SubjectInCurriculum {
			allSubjects = append(allSubjects, GlobalSubject{
				ClassroomID: cl.ID,
				Subject:     sic.SubjectInPreCurriculum,
			})
		}
	}

	// 2. จัดลำดับตามหน่วยกิตจากมากไปน้อย และเรียงตาม SubjectID สำหรับวิชาที่เหมือนกัน
	sort.Slice(allSubjects, func(i, j int) bool {
		if allSubjects[i].Subject.Credit != allSubjects[j].Subject.Credit {
			return allSubjects[i].Subject.Credit > allSubjects[j].Subject.Credit
		}
		if allSubjects[i].Subject.SubjectID != allSubjects[j].Subject.SubjectID {
			return allSubjects[i].Subject.SubjectID < allSubjects[j].Subject.SubjectID
		}
		return allSubjects[i].ClassroomID < allSubjects[j].ClassroomID
	})

	// 3. เริ่มเเจกงานจากกองกลาง
	for _, gs := range allSubjects {
		subject := gs.Subject
		clID := gs.ClassroomID

		if !is_everage {
			if _, ok := classroomAssign[fmt.Sprintf("%d-%d", clID, subject.SubjectID)]; ok {
				continue
			}
		}

		//หาครูที่สอนวิชานี้ได้
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
		//เอาคนคะเเนนเยอะก่อน (Preference 10 -> 1)
		for i := 10; i >= 1; i-- {
			mostPreferanceteacher := []uint{}
			for id, preferance := range teacherTech {
				if preferance == i {
					mostPreferanceteacher = append(mostPreferanceteacher, id)
				}
			}

			if len(mostPreferanceteacher) == 0 {
				continue
			}

			// เรียงตามจำนวนชั่วโมง น้อยไปมาก (ระบุ Priority ใหม่)
			if len(mostPreferanceteacher) > 1 {
				sort.Slice(mostPreferanceteacher, func(i, j int) bool {
					tI := teacherAssign[mostPreferanceteacher[i]]
					tJ := teacherAssign[mostPreferanceteacher[j]]

					// --- เพิ่มระบบ Loyalty Check ---
					// ใครเคยรับวิชานี้ไปแล้ว (ในห้องอื่น) จะมี Priority สูงกว่า
					hasSubI := false
					for _, s := range tI.Subjects {
						if s.SubjectID == subject.SubjectID {
							hasSubI = true
							break
						}
					}
					hasSubJ := false
					for _, s := range tJ.Subjects {
						if s.SubjectID == subject.SubjectID {
							hasSubJ = true
							break
						}
					}

					if hasSubI != hasSubJ {
						return hasSubI // ถ้า I มี แต่ J ไม่มี ให้ I ก่อน
					}

					// ถ้ามีเหมือนกัน หรือไม่มีเหมือนกัน ให้คนงานน้อยก่อนตามปกติ
					return tI.AllocatedHour < tJ.AllocatedHour
				})
			}

			//เเจกงาน
			for _, teacherID := range mostPreferanceteacher {
				if teacherAssign[teacherID].AllocatedHour+subject.Credit > 40 {
					continue
				}

				ta := teacherAssign[teacherID]
				ta.AllocatedHour += subject.Credit
				key := fmt.Sprintf("%d-%d", subject.SubjectID, clID)
				ta.Subjects[key] = SubjectAssign{
					SubjectID:   subject.SubjectID,
					Credit:      subject.Credit,
					ClassroomID: clID,
				}
				teacherAssign[teacherID] = ta
				classroomAssign[fmt.Sprintf("%d-%d", clID, subject.SubjectID)] = struct{}{}
				isBreak = true
				break
			}
			if isBreak {
				break
			}
		}

		if !is_everage {
			if !isBreak {
				return nil, fmt.Errorf("ไม่สามารถจัดตารางได้: วิชาวรหัส %d ในห้องพัก %d ไม่มีครูที่เหมาะสม", subject.SubjectID, clID)
			}
		}
	}

	return &ResultAssign{
		teacherAssign:   teacherAssign,
		classroomAssign: classroomAssign,
		erverageHour:    everageHour,
	}, nil
}

func Can(teachers []entities.Teacher, classrooms []entities.Classroom, rules []TeacherRule) (SchedulRes, error) {
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
	oderTable := orderSchedule(result.teacherAssign, classrooms, rules)
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

func orderSchedule(teacherAssign map[uint]TeacherAssign, classrooms []entities.Classroom, rules []TeacherRule) map[position]*TeacherInOrder {
	teacher := NewTeacherAllocate(teacherAssign, rules)
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
	TeacherAssign      map[uint]TeacherAssign
	YAdded             map[int]InAdded
	TeacherRules       map[uint]TeacherRule
	LastClassroomAlloc map[uint]LastAlloc // [ClassroomID]LastAlloc
}

type LastAlloc struct {
	TeacherID uint
	SubjectID uint
	Count     int
}
type InAdded struct {
	teachers map[uint]struct{}
}

func NewTeacherAllocate(TeacherAssign map[uint]TeacherAssign, rules []TeacherRule) teacherAllocate {
	yAdded := make(map[int]InAdded)
	for i := 1; i <= 40; i++ {
		teachers := make(map[uint]struct{})
		for teacherID := range TeacherAssign {
			teachers[teacherID] = struct{}{}
		}
		yAdded[i] = InAdded{teachers: teachers}
	}

	ruleMap := make(map[uint]TeacherRule)
	for _, r := range rules {
		existing := ruleMap[r.TeacherID]
		if r.Consecutive > 0 {
			existing.Consecutive = r.Consecutive
		}
		existing.BusyPeriods = append(existing.BusyPeriods, r.BusyPeriods...)
		existing.PreferredPeriods = append(existing.PreferredPeriods, r.PreferredPeriods...)
		existing.TeacherID = r.TeacherID
		ruleMap[r.TeacherID] = existing
	}

	return teacherAllocate{
		TeacherAssign:      TeacherAssign,
		YAdded:             yAdded,
		TeacherRules:       ruleMap,
		LastClassroomAlloc: make(map[uint]LastAlloc),
	}
}

func (ta *teacherAllocate) tik(y int, classroomID uint) *TeacherInOrder {
	var winID uint
	var winSubjectID uint
	var winKey string
	var maxScore int = -10000
	found := false

	last := ta.LastClassroomAlloc[classroomID]

	// 1. ลองหาจากครูที่ถูก "มอบหมาย" งานไว้ก่อน (Primary Assigned)
	if len(ta.YAdded[y].teachers) > 0 {
		for teacherID := range ta.YAdded[y].teachers {
			t := ta.TeacherAssign[teacherID]
			rule := ta.TeacherRules[teacherID]

			isBusy := false
			for _, p := range rule.BusyPeriods {
				if p == y {
					isBusy = true
					break
				}
			}

			// ค้นหาวิจาที่ "ดีที่สุด" (ID น้อยสุด) ของครูคนนี้ในห้องนี้
			var bestSub SubjectAssign
			var bestKey string
			minSubID := uint(999999)
			for key, sub := range t.Subjects {
				if sub.ClassroomID == classroomID && sub.Credit > 0 {
					if sub.SubjectID < minSubID {
						minSubID = sub.SubjectID
						bestSub = sub
						bestKey = key
					}
				}
			}

			if minSubID != 999999 {
				sub := bestSub
				key := bestKey

				// คะแนนพื้นฐานใช้ AllocatedHour
				currentScore := t.AllocatedHour * 5

				if isBusy {
					currentScore -= 10000
				}

				// กฎเหล็ก: ถ้าเป็นครูคนเดิมในห้องเดิม แต่เปลี่ยนวิชา ให้โดนหักคะแนนหนัก
				// (ยกเว้นวิชาเดิมจะหมดแล้วจริงๆ ซึ่งเช็คจากการที่ ID เปลี่ยนไป)
				if last.TeacherID == teacherID && last.SubjectID != sub.SubjectID {
					currentScore -= 5000 // ป้องกันการสลับวิชาไปมา
				}

				// โบนัสความชอบ (Preferred Periods)
				for _, pp := range rule.PreferredPeriods {
					if pp == y {
						currentScore += 500
						break
					}
				}

				// โบนัสความต่อเนื่อง (เฉพาะวิชาเดิม)
				if last.TeacherID == teacherID && last.SubjectID == sub.SubjectID {
					target := 2
					if rule.Consecutive > 0 {
						target = rule.Consecutive
					}
					if last.Count < target {
						currentScore += 2000
					}
				}

				if currentScore > maxScore {
					maxScore = currentScore
					winID = teacherID
					winSubjectID = sub.SubjectID
					winKey = key
					found = true
				}
			}
		}
	}

	// 2. [Substitution] ถ้าคนหลักติด Busy หรือยังหาคนลงไม่ได้ ให้ลองคนอื่นที่ "สอนได้" และ "ว่างจริง"
	if !found || maxScore < 0 {
		subScore := -10000
		var subID uint
		var subTID uint
		subFound := false

		neededSubjects := make(map[uint]int)
		for _, t := range ta.TeacherAssign {
			for _, sub := range t.Subjects {
				if sub.ClassroomID == classroomID && sub.Credit > 0 {
					neededSubjects[sub.SubjectID] += sub.Credit
				}
			}
		}

		for sID, credit := range neededSubjects {
			for tID := range ta.YAdded[y].teachers {
				t := ta.TeacherAssign[tID]
				rule := ta.TeacherRules[tID]

				busy := false
				for _, p := range rule.BusyPeriods {
					if p == y {
						busy = true
						break
					}
				}
				if busy {
					continue
				}

				if _, canTeach := t.MySubjects[sID]; canTeach {
					if t.AllocatedHour > 0 {
						// คะแนนตัวสำรอง: เน้นกระจายงาน
						currentScore := (t.AllocatedHour * 3) + (credit * 10)

						if currentScore > subScore {
							subScore = currentScore
							subTID = tID
							subID = sID
							subFound = true
						}
					}
				}
			}
		}

		// ถ้าตัวสำรองคะแนนดีกว่า (และไม่ติด Busy) ให้เลือกตัวสำรอง
		if subFound && subScore > maxScore {
			maxScore = subScore
			winID = subTID
			winSubjectID = subID
			winKey = ""
			found = true
		}
	}

	// ถ้าเจอกันแต่ติด Busy มหาศาล ให้ปล่อยว่าง
	if !found || maxScore < -5000 {
		ta.LastClassroomAlloc[classroomID] = LastAlloc{}
		return nil
	}

	// อัปเดตข้อมูลการสอน
	if last.TeacherID == winID && last.SubjectID == winSubjectID {
		last.Count++
	} else {
		last = LastAlloc{TeacherID: winID, SubjectID: winSubjectID, Count: 1}
	}
	ta.LastClassroomAlloc[classroomID] = last

	delete(ta.YAdded[y].teachers, winID)
	teacher := ta.TeacherAssign[winID]
	teacher.AllocatedHour -= 1

	// หักโควตาวิจา (ถ้าเป็นการสอนแทน ต้องไปหาฟิลด์หักให้ถูกตัว)
	if winKey != "" {
		if s, ok := teacher.Subjects[winKey]; ok {
			s.Credit -= 1
			teacher.Subjects[winKey] = s
		}
	} else {
		// กรณีตัวสำรอง: หัก Credit จากครูคนแรกที่ได้รับมอบหมายวิชานี้ในห้องนี้ (เพื่อยอดรวมไม่เพี้ยน)
		for tID, t := range ta.TeacherAssign {
			for k, s := range t.Subjects {
				if s.ClassroomID == classroomID && s.SubjectID == winSubjectID && s.Credit > 0 {
					s.Credit -= 1
					ta.TeacherAssign[tID].Subjects[k] = s
					goto updated
				}
			}
		}
	}

updated:
	ta.TeacherAssign[winID] = teacher
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
