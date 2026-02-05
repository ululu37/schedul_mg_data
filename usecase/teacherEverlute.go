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

		subjectsEvu := allSubjects

		for _, subject := range mysubject {
			for j, s := range subjectsEvu {
				if subject.SubjectID == s.ID {
					subjectsEvu = append(subjectsEvu[:j], subjectsEvu[j+1:]...)
					break
				}
			}
		}

		for p := 1; ; p++ {

			limit := 100
			for len(subjectsEvu) > 0 {
				n := limit
				if len(subjectsEvu) < limit {
					n = len(subjectsEvu)
				}
				t.everluteAi(teacher, subjectsEvu[:n])
				subjectsEvu = subjectsEvu[n:]
			}
			t.everluteAi(teacher, subjectsEvu[:100])
			subjectsEvu = subjectsEvu[100:]

			if len(subjectsEvu) <= 100 {
				if len(subjectsEvu) > 0 {
					t.everluteAi(teacher, subjectsEvu)
				}

				break
			}

		}

	}

	return nil
}

func (t *TeacherEverlute) everluteAi(teacher entities.Teacher, mysubject []entities.Subject) (*dto.EvaluationResponse, error) {

	respBody, errAi := t.agent.Chat([]aiAgent.Message{
		{
			Role: "system",
			Content: `
You are an aptitude evaluator for teachers and subjects.

Evaluate the teacher’s aptitude for each subject.

Aptitude score scale: 1–10

มาตราฐานการประเมิน
 10 - 9 > บอกว่าถนัด
 8-6 > ควรสอนได้ตามสายงานที่จบมา
 5-1 > น่าจะสอนได้ 
 วิชาที่ไม่น่าจะสอนไม่ได้ ไห้ตัดออกไม่ต้องส่งมา

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

	if errAi != nil {
		return nil, errAi
	}
	var res dto.EvaluationResponse

	errJsonEncode := json.Unmarshal([]byte(respBody.Choices[0].Message.Content), &res)
	if errJsonEncode != nil {
		return nil, errJsonEncode
	}
	for _, ev := range res.Evaluation {
		fmt.Printf("subject_id=%d aptitude=%d\n", ev.ID, ev.Adtritud)
	}
	fmt.Printf("-----------")
	return &res, nil
}
