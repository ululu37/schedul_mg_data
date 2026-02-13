// NewTeacherEverlute creates a new TeacherEverlute instance

package usecase

import (
	"encoding/json"
	"fmt"
	dto "scadulDataMono/domain/DTO"
	"scadulDataMono/domain/entities"
	aiAgent "scadulDataMono/infra/Agent"
	"strings"
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

			MysDB, _, err := t.teacherMg.GetMySubject(teacher.ID, -2, "", 1, 100)
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
			//fmt.Println("len subjectsEvuList before: %d, n : %d", len(subjectsEvuList), n)
			aiRes, err := t.everluteAi(teacher, subjectsEvuList[:n])
			if err != nil {
				fmt.Println("errAI", err)
				return err
			}
			for _, ev := range aiRes.Evaluation {
				fmt.Printf("id: %v, aptitude: %v\n", ev.ID, ev.Aptitude)
				toMySubject = append(toMySubject, entities.TeacherMySubject{
					TeacherID:  teacher.ID,
					SubjectID:  uint(ev.ID),
					Preference: ev.Aptitude,
				})

			}

			subjectsEvuList = subjectsEvuList[n:]
		}
		fmt.Printf("toMySubject %+v\n", toMySubject)
		t.teacherMg.AddMySubject(teacher.ID, toMySubject)

	}
	fmt.Println("success")
	return nil
}

func (t *TeacherEverlute) everluteAi(teacher entities.Teacher, mysubject []entities.Subject) (*dto.EvaluationResponse, error) {
	fmt.Println("everluteAi teacher: ", teacher.Name, "Resume: ", teacher.Resume)
	//fmt.Println("everluteAi mysubject: ", mysubject)
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
    { "id": "number", "aptitude": "number" }
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
	if respBody == nil || len(respBody.Choices) == 0 {
		return nil, fmt.Errorf("AI agent returned no choices")
	}

	var res dto.EvaluationResponse

	content := respBody.Choices[0].Message.Content
	start := strings.Index(content, "{")
	end := strings.LastIndex(content, "}")
	if start != -1 && end != -1 && start < end {
		content = content[start : end+1]
	}

	errJsonEncode := json.Unmarshal([]byte(content), &res)
	if errJsonEncode != nil {
		return nil, errJsonEncode
	}
	fmt.Println("AIres")
	return &res, nil
}
