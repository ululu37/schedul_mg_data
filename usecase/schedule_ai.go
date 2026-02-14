package usecase

import (
	"encoding/json"
	"fmt"
	"scadulDataMono/domain/entities"
	aiAgent "scadulDataMono/infra/Agent"
	"strings"
)

type ScheduleAI struct {
	agent *aiAgent.Agent
}

func NewScheduleAI(agent *aiAgent.Agent) *ScheduleAI {
	return &ScheduleAI{agent: agent}
}

func (s *ScheduleAI) ParsePreferences(userInput string, teachers []entities.Teacher) ([]TeacherRule, error) {
	allRules := []TeacherRule{}

	batchSize := 20
	for i := 0; i < len(teachers); i += batchSize {
		end := i + batchSize
		if end > len(teachers) {
			end = len(teachers)
		}
		batch := teachers[i:end]

		teacherList := ""
		for _, t := range batch {
			teacherList += fmt.Sprintf("- ID: %d, Name: %s\n", t.ID, t.Name)
		}

		prompt := fmt.Sprintf(`
You are a scheduling assistant. Convert the user's request into specific rules for the teachers in the PROVIDED LIST below.

Strict ID Mapping Rules:
1. **ONLY** create rules for Teacher IDs that appear in the list below.
2. If a teacher mentioned in the request (e.g., "%s") is **NOT** in the list below, ignore that teacher for this response. Do not guess IDs.
3. Use "teacher_id": 0 **ONLY** if the request is a global command (e.g., "all teachers", "everyone").

Request Logic:
- "consecutive": Number of hours in a row.
- "busy_periods": Periods they CANNOT teach (1-40).
- "preferred_periods": Periods they WANT to teach (1-40).
   - Morning: 1-4 (Mon), 9-12 (Tue), 17-20 (Wed), 25-28 (Thu), 33-36 (Fri).
   - Afternoon: 5-8 (Mon), 13-16 (Tue), 21-24 (Wed), 29-32 (Thu), 37-40 (Fri).

Periods Context:
1-8 Mon, 9-16 Tue, 17-24 Wed, 25-32 Thu, 33-40 Fri.

Provided Teacher List in this Batch:
%s

User Request: "%s"

Output ONLY a JSON array:
[
  {"teacher_id": ID, "consecutive": 0, "busy_periods": [], "preferred_periods": []}
]
`, userInput, teacherList, userInput)

		resp, err := s.agent.Chat([]aiAgent.Message{
			{Role: "system", Content: "Output valid JSON only. NEVER map a name to an ID unless that ID is in the provided list. If no relevant teachers in list, return []."},
			{Role: "user", Content: prompt},
		})
		if err != nil {
			continue
		}
		if resp == nil || len(resp.Choices) == 0 {
			continue
		}

		content, ok := resp.Choices[0].Message.Content.(string)
		if !ok {
			continue
		}

		startPos := strings.Index(content, "[")
		endPos := strings.LastIndex(content, "]")
		if startPos == -1 || endPos == -1 || startPos >= endPos {
			continue
		}
		content = content[startPos : endPos+1]

		var results []struct {
			TeacherID        uint  `json:"teacher_id"`
			Consecutive      int   `json:"consecutive"`
			BusyPeriods      []int `json:"busy_periods"`
			PreferredPeriods []int `json:"preferred_periods"`
		}

		if err := json.Unmarshal([]byte(content), &results); err != nil {
			continue
		}

		for _, r := range results {
			if r.TeacherID == 0 {
				for _, t := range batch {
					allRules = append(allRules, TeacherRule{
						TeacherID:        t.ID,
						Consecutive:      r.Consecutive,
						BusyPeriods:      r.BusyPeriods,
						PreferredPeriods: r.PreferredPeriods,
					})
				}
			} else {
				// Verify if ID belongs to this batch
				belongs := false
				for _, t := range batch {
					if t.ID == r.TeacherID {
						belongs = true
						break
					}
				}
				if belongs {
					allRules = append(allRules, TeacherRule{
						TeacherID:        r.TeacherID,
						Consecutive:      r.Consecutive,
						BusyPeriods:      r.BusyPeriods,
						PreferredPeriods: r.PreferredPeriods,
					})
				}
			}
		}
	}
	fmt.Println("allRules", allRules)
	return allRules, nil
}
