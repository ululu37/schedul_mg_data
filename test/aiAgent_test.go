package gpt_test

import (
	"encoding/json"
	"fmt"
	"log"
	dto "scadulDataMono/domain/DTO"
	aiAgent "scadulDataMono/infra/Agent"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGpt4oMini_Chat(t *testing.T) {
	client := aiAgent.Agent{
		ApiURL: "https://openrouter.ai/api/v1/chat/completions",
		ApiKey: "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",
	}

	respBody, err := client.Chat([]aiAgent.Message{
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
			Role: "user",
			Content: `{
		    "teacher_resume": ["ครูคอม"],
		    "subjects": [
  { "id": 1, "name": "data structures" },
  { "id": 2, "name": "english" },
  { "id": 3, "name": "การตลาด" },
  { "id": 4, "name": "computer networks" },
  { "id": 5, "name": "วิทยาศาสตร์" },
  { "id": 6, "name": "database systems" },
  { "id": 7, "name": "algorithms" },
  { "id": 8, "name": "software engineering" },
  { "id": 9, "name": "web development" },
  { "id": 10, "name": "cloud computing" },

  { "id": 11, "name": "cyber security" },
  { "id": 12, "name": "computer architecture" },
  { "id": 13, "name": "artificial intelligence" },
  { "id": 14, "name": "machine learning" },
  { "id": 15, "name": "data science" },
  { "id": 16, "name": "statistics" },
  { "id": 17, "name": "linear algebra" },
  { "id": 18, "name": "calculus" },
  { "id": 19, "name": "discrete mathematics" },
  { "id": 20, "name": "probability" },

  { "id": 21, "name": "เศรษฐศาสตร์" },
  { "id": 22, "name": "บัญชีพื้นฐาน" },
  { "id": 23, "name": "การเงิน" },
  { "id": 24, "name": "การบริหารธุรกิจ" },
  { "id": 25, "name": "ผู้ประกอบการ" },
  { "id": 26, "name": "การจัดการโครงการ" },
  { "id": 27, "name": "การวิเคราะห์ธุรกิจ" },
  { "id": 28, "name": "digital marketing" },
  { "id": 29, "name": "content marketing" },
  { "id": 30, "name": "branding" },

  { "id": 31, "name": "user experience design" },
  { "id": 32, "name": "user interface design" },
  { "id": 33, "name": "graphic design" },
  { "id": 34, "name": "multimedia" },
  { "id": 35, "name": "animation" },
  { "id": 36, "name": "game development" },
  { "id": 37, "name": "game design" },
  { "id": 38, "name": "computer graphics" },
  { "id": 39, "name": "3d modeling" },
  { "id": 40, "name": "virtual reality" },

  { "id": 41, "name": "augmented reality" },
  { "id": 42, "name": "mobile application development" },
  { "id": 43, "name": "ios development" },
  { "id": 44, "name": "android development" },
  { "id": 45, "name": "cross platform development" },
  { "id": 46, "name": "devops" },
  { "id": 47, "name": "docker" },
  { "id": 48, "name": "kubernetes" },
  { "id": 49, "name": "site reliability engineering" },
  { "id": 50, "name": "system administration" },

  { "id": 51, "name": "network administration" },
  { "id": 52, "name": "linux systems" },
  { "id": 53, "name": "windows systems" },
  { "id": 54, "name": "server management" },
  { "id": 55, "name": "distributed systems" },
  { "id": 56, "name": "microservices" },
  { "id": 57, "name": "api design" },
  { "id": 58, "name": "restful services" },
  { "id": 59, "name": "grpc" },
  { "id": 60, "name": "message queues" },

  { "id": 61, "name": "software testing" },
  { "id": 62, "name": "unit testing" },
  { "id": 63, "name": "integration testing" },
  { "id": 64, "name": "test automation" },
  { "id": 65, "name": "quality assurance" },
  { "id": 66, "name": "software architecture" },
  { "id": 67, "name": "design patterns" },
  { "id": 68, "name": "clean code" },
  { "id": 69, "name": "refactoring" },
  { "id": 70, "name": "code review" },

  { "id": 71, "name": "technical writing" },
  { "id": 72, "name": "presentation skills" },
  { "id": 73, "name": "communication skills" },
  { "id": 74, "name": "critical thinking" },
  { "id": 75, "name": "problem solving" },
  { "id": 76, "name": "leadership" },
  { "id": 77, "name": "team management" },
  { "id": 78, "name": "agile methodology" },
  { "id": 79, "name": "scrum" },
  { "id": 80, "name": "kanban" },

  { "id": 81, "name": "blockchain" },
  { "id": 82, "name": "smart contracts" },
  { "id": 83, "name": "fintech" },
  { "id": 84, "name": "cryptography" },
  { "id": 85, "name": "information security" },
  { "id": 86, "name": "ethical hacking" },
  { "id": 87, "name": "penetration testing" },
  { "id": 88, "name": "risk management" },
  { "id": 89, "name": "compliance" },
  { "id": 90, "name": "it governance" },

  { "id": 91, "name": "internet of things" },
  { "id": 92, "name": "embedded systems" },
  { "id": 93, "name": "robotics" },
  { "id": 94, "name": "automation systems" },
  { "id": 95, "name": "industrial technology" },
  { "id": 96, "name": "smart factory" },
  { "id": 97, "name": "supply chain management" },
  { "id": 98, "name": "logistics" },
  { "id": 99, "name": "innovation management" },
  { "id": 100, "name": "technology strategy" }
		    ] 
		}`,
		},
	})

	var res dto.EvaluationResponse

	errJsonEncode := json.Unmarshal([]byte(respBody.Choices[0].Message.Content), &res)
	if errJsonEncode != nil {
		log.Fatal(errJsonEncode)
	}

	// ใช้งานต่อ
	for _, ev := range res.Evaluation {
		fmt.Printf("subject_id=%d aptitude=%d\n", ev.ID, ev.Adtritud)
	}

	//fmt.Printf("%s", respBody.Choices[0].Message.Content)
	assert.NoError(t, err)
	assert.NotNil(t, respBody)
	assert.NotEmpty(t, respBody.Choices)
	assert.NotEmpty(t, respBody.Choices[0].Message.Content)
}
