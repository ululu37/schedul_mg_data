// NewTeacherEverlute creates a new TeacherEverlute instance

package usecase

import (
	"encoding/json"
	"fmt"
	"regexp"
	dto "scadulDataMono/domain/DTO"
	"scadulDataMono/domain/entities"
	aiAgent "scadulDataMono/infra/Agent"
	"strings"
	"sync" // Added for concurrency
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
func (t *TeacherEverlute) Everlute(progress chan string) error {
	report := func(msg string) {
		if progress != nil {
			progress <- msg
		}
		fmt.Println(msg)
	}

	report("Fetching teacher list...")
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

	report(fmt.Sprintf("Found %d teachers. Fetching subjects...", len(teachers)))
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

	report(fmt.Sprintf("Evaluating %d teachers across %d subjects...", len(teachers), len(allSubjects)))
	// Concurrency control: limit to 10 concurrent AI requests
	sem := make(chan struct{}, 10)
	var wg sync.WaitGroup

	for _, teacher := range teachers {
		wg.Add(1)
		go func(tchr entities.Teacher) {
			defer wg.Done()

			// get mysubject
			mysubject := []entities.TeacherMySubject{}
			for p := 1; ; p++ {
				MysDB, _, err := t.teacherMg.GetMySubject(tchr.ID, -2, "", p, 100)
				if err != nil {
					report(fmt.Sprintf("err getting subjects for teacher %s: %v", tchr.Name, err))
					return
				}
				if len(MysDB) == 0 {
					break
				}
				mysubject = append(mysubject, MysDB...)
				if len(MysDB) < 100 {
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

			if len(subjectsEvuList) == 0 {
				return
			}

			report(fmt.Sprintf("Teacher %s: evaluating %d remaining subjects", tchr.Name, len(subjectsEvuList)))

			var batchWg sync.WaitGroup
			limit := 100
			for len(subjectsEvuList) > 0 {
				n := limit
				if len(subjectsEvuList) < limit {
					n = len(subjectsEvuList)
				}

				batch := subjectsEvuList[:n]
				subjectsEvuList = subjectsEvuList[n:]

				batchWg.Add(1)
				go func(innerTchr entities.Teacher, bth []entities.Subject) {
					defer batchWg.Done()

					sem <- struct{}{}        // Acquire semaphore
					defer func() { <-sem }() // Release semaphore

					aiRes, err := t.everluteAi(innerTchr, bth)
					if err != nil {
						report(fmt.Sprintf("errAI for teacher %s: %v", innerTchr.Name, err))
						return
					}

					batchResults := []entities.TeacherMySubject{}
					for _, ev := range aiRes.Evaluation {
						// report(fmt.Sprintf("teacher: %s, id: %v, aptitude: %v", innerTchr.Name, ev.ID, ev.Aptitude))
						batchResults = append(batchResults, entities.TeacherMySubject{
							TeacherID:  innerTchr.ID,
							SubjectID:  uint(ev.ID),
							Preference: ev.Aptitude,
						})
					}

					if len(batchResults) > 0 {
						if err := t.teacherMg.AddMySubject(innerTchr.ID, batchResults); err != nil {
							report(fmt.Sprintf("err saving results for teacher %s: %v", innerTchr.Name, err))
						}
					}
				}(tchr, batch)
			}
			batchWg.Wait()
			report(fmt.Sprintf("Teacher %s: evaluation complete", tchr.Name))
		}(teacher)
	}
	wg.Wait()
	report("AI Evaluation Process Complete")
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
You are a teacher aptitude evaluator. Output ONLY valid JSON.
Evaluate the teacher's aptitude for each provided subject based on their resume.
Aptitude score scale: 0–10

Scoring Guide:
- 10-9: Expert/Fluent in this subject.
- 8-6: Should be able to teach based on their field of study.
- 5-1: Might be able to teach with some preparation.
- 0: Unfit to teach this subject (MUST include these in the output too).

Output EXACTLY this JSON structure:
{
  "evaluation": [
    { "id": number, "aptitude": number }
  ]
}

IMPORTANT: Do NOT include any explanations, Markdown formatting (e.g., no json blocks), or trailing commas.
`,
		},
		{
			Role:    "user",
			Content: fmt.Sprintf("teacher_resume: %s\nsubjects_to_evaluate: %+v", teacher.Resume, mysubject),
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

	// Handle content type assertion
	content, ok := respBody.Choices[0].Message.Content.(string)
	if !ok {
		return nil, fmt.Errorf("AI agent returned non-string content")
	}

	// Clean up potential Markdown blocks
	content = strings.TrimSpace(content)
	if strings.HasPrefix(content, "```json") {
		content = strings.TrimPrefix(content, "```json")
		content = strings.TrimSuffix(content, "```")
	} else if strings.HasPrefix(content, "```") {
		content = strings.TrimPrefix(content, "```")
		content = strings.TrimSuffix(content, "```")
	}

	start := strings.Index(content, "{")
	end := strings.LastIndex(content, "}")
	if start != -1 && end != -1 && start < end {
		content = content[start : end+1]
	}

	// Remove trailing commas that AI often generates (e.g., [1, 2, ])
	re := regexp.MustCompile(`,\s*([\]}])`)
	content = re.ReplaceAllString(content, "$1")

	errJsonEncode := json.Unmarshal([]byte(content), &res)
	if errJsonEncode != nil {
		fmt.Printf("DEBUG: Failed to parse AI JSON for teacher %s. Raw content: %s\n", teacher.Name, content)
		return nil, fmt.Errorf("failed to parse AI response: %v", errJsonEncode)
	}
	fmt.Println("AI evaluation parsed successfully")
	return &res, nil
}
