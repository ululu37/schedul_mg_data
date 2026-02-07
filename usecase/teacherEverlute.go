// NewTeacherEverlute creates a new TeacherEverlute instance

package usecase

import (
	"encoding/json"
	"fmt"
	dto "scadulDataMono/domain/DTO"
	"scadulDataMono/domain/entities"
	aiAgent "scadulDataMono/infra/Agent"
)

type TeacherEverlute struct {
	teacherMg *TeacherMg
	subjectMg *SubjectMg
	agent     *aiAgent.Agent
}

func NewTeacherEverlute(teacherMg *TeacherMg, subjectMg *SubjectMg, agent *aiAgent.Agent) *TeacherEverlute {
	return &TeacherEverlute{
		teacherMg: teacherMg,
		subjectMg: subjectMg,
		agent:     agent,
	}
}
func (t *TeacherEverlute) Everlute() error {
	// get teacher list
	teachers := []entities.Teacher{}
	for p := 1; ; p++ {
		teachersDB, _, err := t.teacherMg.Listing("", p, 100)
		if err != nil {
			return err
		}
		if len(teachersDB) == 0 {
			break
		}
		teachers = append(teachers, teachersDB...)
		if len(teachersDB) < 100 {
			break
		}
	}

	allSubjects := []entities.Subject{}
	// Get all subjects
	for p := 1; ; p++ {
		subjects, _, sErr := t.subjectMg.Listing("", p, 100)

		if sErr != nil {
			return sErr
		}

		allSubjects = append(allSubjects, subjects...)
		if len(subjects) < 100 {
			break
		}
	}
	//fmt.Println("all::", allSubjects)
	// Evaluate each teacher

	for _, teacher := range teachers {

		// get mysubject
		mysubject := []entities.TeacherMySubject{}
		for p := 1; ; p++ {

			MysDB, _, err := t.teacherMg.GetMySubject(teacher.ID, -2, 1, 100)
			if len(MysDB) == 0 {
				break
			}
			if err != nil {
				return err
			}
			mysubject = append(mysubject, MysDB...)
			if len(MysDB) <= 100 {
				break
			}
		}

		subjectsEvu := make(map[uint]entities.Subject)

		for _, s := range allSubjects {
			subjectsEvu[s.ID] = s
		}

		for _, subject := range mysubject {
			delete(subjectsEvu, subject.SubjectID)
		}
		subjectsEvuList := []entities.Subject{}
		for _, v := range subjectsEvu {
			subjectsEvuList = append(subjectsEvuList, v)
		}

		toMySubject := []entities.TeacherMySubject{}
		//	fmt.Println(":::::::", subjectsEvu)
		limit := 100
		for len(subjectsEvuList) > 0 {
			n := limit

			if len(subjectsEvuList) < limit {
				n = len(subjectsEvuList)
			}

			if len(subjectsEvuList) == 0 {
				break
			}
			fmt.Println("len subjectsEvuList before: %d, n : %d", len(subjectsEvuList), n)
			aiRes, err := t.everluteAi(teacher, subjectsEvuList[:n])
			if err != nil {
				return err
			}
			for _, ev := range aiRes.Evaluation {
				fmt.Println("id: %v, adtitude: %v", ev.ID, ev.Adtritud)
				toMySubject = append(toMySubject, entities.TeacherMySubject{
					SubjectID:  uint(ev.ID),
					Preference: ev.Adtritud,
				})

			}

			subjectsEvuList = subjectsEvuList[n:]
		}

		t.teacherMg.AddMySubject(teacher.ID, toMySubject)
	}

	return nil
}

func (t *TeacherEverlute) everluteAi(teacher entities.Teacher, mysubject []entities.Subject) (*dto.EvaluationResponse, error) {
	//fmt.Printf("ssssssss\n", mysubject)
	respBody, errAi := t.agent.Chat([]aiAgent.Message{
		{
			Role: "system",
			Content: `
You are an aptitude evaluator for teachers and subjects.

Evaluate the teacher’s aptitude for each subject.

Aptitude score scale: 0–10

มาตราฐานการประเมิน
 10 - 9 > บอกว่าถนัด
 8-6 > ควรสอนได้ตามสายงานที่จบมา
 5-1 > น่าจะสอนได้ 
 วิชาที่ไม่น่าจะสอนไม่ได้ ไห้คะเเนเป็น 0 เเล้วตส่งมาด้วย
Output JSON only.

Output schema:
{
  "evaluation": [
    { "id": "number", "adtritud": "number" }
  ]
}
`,
		},
		{
			Role:    "user",
			Content: fmt.Sprintf("teacher_resume: %s,%+v", teacher.Resume, mysubject),
		},
	},
	)
	//fmt.Printf("rrrrrrrrr%+v\n", respBody)
	if errAi != nil {
		return nil, errAi
	}

	var res dto.EvaluationResponse

	errJsonEncode := json.Unmarshal([]byte(respBody.Choices[0].Message.Content), &res)
	if errJsonEncode != nil {
		return nil, errJsonEncode
	}
	// for _, ev := range res.Evaluation {
	// 	fmt.Printf("subject_id=%d aptitude=%d\n", ev.ID, ev.Adtritud)
	// }
	fmt.Printf("-----------")
	return &res, nil
}
