package usecase

import (
	"encoding/json"
	"fmt"
	"scadulDataMono/domain/entities"
	aiAgent "scadulDataMono/infra/Agent"
	"strings"
)

type ImportPrecuriculum struct {
	preCurriculum *PreCurriculum
	agent         *aiAgent.Agent
}

func NewImportPrecuriculum(preCurriculum *PreCurriculum, agent *aiAgent.Agent) *ImportPrecuriculum {
	return &ImportPrecuriculum{
		preCurriculum: preCurriculum,
		agent:         agent,
	}
}

type ImportResult struct {
	Name     string `json:"name"`
	Subjects []struct {
		Name   string `json:"name"`
		Credit int    `json:"credit"`
	} `json:"subjects"`
}

func (u *ImportPrecuriculum) Import(input interface{}, progress chan string) (uint, error) {
	report := func(msg string) {
		if progress != nil {
			progress <- msg
		}
		fmt.Println(msg)
	}

	report("Sending curriculum data to AI for extraction...")
	prompt := `
Extract ALL curriculum information from the provided input (text or image) and return it only in JSON format.
The input may contain multiple pages or long tables. DO NOT skip any subjects.

JSON schema:
{
  "name": "Full Name of the Curriculum/Program",
  "subjects": [
    { "name": "Exact Subject Name", "credit": number }
  ]
}

CRITICAL RULES:
1. 'credit' MUST be the 'Teaching Hours per Week' (ชั่วโมงการสอนต่อสัปดาห์), usually found in columns labeled as (ท-ป-ศ) or (Lecture-Lab-SelfStudy). It is NOT the total academic credits.
2. If you see (3-0-6), the credit is 3. If you see (1-2-3), the credit is 3 (1+2).
3. If 'credit' is not clearly visible for a subject, estimate it based on similar subjects (usually 2 or 3).
4. Extract EVERY subject listed in the document. Do not summarize or omit any.
5. If the input is from a PDF or long text, scan through the entire content.
6. Return ONLY the JSON object.
`

	resp, err := u.agent.Chat([]aiAgent.Message{
		{Role: "system", Content: "You are a curriculum data extractor. Output JSON only."},
		{
			Role: "user",
			Content: []interface{}{
				map[string]interface{}{
					"type": "text",
					"text": prompt,
				},
				input,
			},
		},
	})
	if err != nil {
		fmt.Printf("AI agent error: %v\n", err) // เพิ่มบรรทัดนี้เพื่อดูสาเหตุ
		return 0, err
	}
	if resp == nil || len(resp.Choices) == 0 {
		fmt.Println("AI agent returned no choices")
		return 0, fmt.Errorf("AI agent returned no choices")
	}

	report("AI successfully extracted data. Parsing results...")
	rawContent := resp.Choices[0].Message.Content
	content, ok := rawContent.(string)
	if !ok {
		return 0, fmt.Errorf("AI response is not a string")
	}

	// Extract JSON
	start := strings.Index(content, "{")
	end := strings.LastIndex(content, "}")
	if start == -1 || end == -1 || start >= end {
		return 0, fmt.Errorf("invalid AI response: %s", content)
	}
	content = content[start : end+1]

	var res ImportResult
	if err := json.Unmarshal([]byte(content), &res); err != nil {
		return 0, err
	}

	report(fmt.Sprintf("Importing curriculum: %s with %d subjects", res.Name, len(res.Subjects)))
	// Save to DB
	preID, err := u.preCurriculum.Create(res.Name)
	if err != nil {
		return 0, err
	}

	newSubjectInCurriculum := []entities.SubjectInPreCurriculum{}
	for _, s := range res.Subjects {
		newSubjectInCurriculum = append(newSubjectInCurriculum, entities.SubjectInPreCurriculum{
			PreCurriculumID: preID,
			Subject:         entities.Subject{Name: s.Name},
			Credit:          s.Credit,
		})
	}

	if err := u.preCurriculum.CreateSubject(preID, newSubjectInCurriculum); err != nil {
		return 0, err
	}

	report("Import Complete")
	return preID, nil
}
