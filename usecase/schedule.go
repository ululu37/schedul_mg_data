package usecase

import (
	"fmt"
	"scadulDataMono/domain/entities"
	"scadulDataMono/infra/gormDB/repo"
	"strings"
)

type ScheduleUsecase struct {
	TeacherRepo        *repo.TeacherRepo
	ClassroomRepo      *repo.ClassroomRepo
	ScadulStudentRepo  *repo.ScadulStudentRepo
	ScadulTeacherRepo  *repo.ScadulTeacherRepo
	ScheduleAI         *ScheduleAI
	SchedulePromptRepo *repo.SchedulePromptRepo
}

func NewScheduleUsecase(
	teacherRepo *repo.TeacherRepo,
	classroomRepo *repo.ClassroomRepo,
	scadulStudentRepo *repo.ScadulStudentRepo,
	scadulTeacherRepo *repo.ScadulTeacherRepo,
	scheduleAI *ScheduleAI,
	schedulePromptRepo *repo.SchedulePromptRepo,
) *ScheduleUsecase {
	return &ScheduleUsecase{
		TeacherRepo:        teacherRepo,
		ClassroomRepo:      classroomRepo,
		ScadulStudentRepo:  scadulStudentRepo,
		ScadulTeacherRepo:  scadulTeacherRepo,
		ScheduleAI:         scheduleAI,
		SchedulePromptRepo: schedulePromptRepo,
	}
}

func (u *ScheduleUsecase) GenerateScheduleWithAI(userInput string) (*SchedulRes, error) {
	// 0. Clear old schedules first
	if err := u.ScadulStudentRepo.DeleteAll(); err != nil {
		return nil, err
	}
	if err := u.ScadulTeacherRepo.DeleteAll(); err != nil {
		return nil, err
	}

	// 1. Save and Fetch Prompts (Cumulative Context)
	if userInput != "" {
		if strings.Contains(userInput, "ล้างคำสั่ง") || strings.Contains(userInput, "clear") {
			u.SchedulePromptRepo.DeleteAll()
			userInput = ""
		} else {
			u.SchedulePromptRepo.Create(&entities.SchedulePrompt{Prompt: userInput})
		}
	}

	allPrompts, _ := u.SchedulePromptRepo.GetAll()
	combinedPrompt := ""
	for _, p := range allPrompts {
		combinedPrompt += p.Prompt + " "
	}

	// 2. Load Data
	teachers, err := u.TeacherRepo.GetAllWithRelations()
	if err != nil {
		return nil, err
	}

	classrooms, err := u.ClassroomRepo.GetAllWithRelations()
	if err != nil {
		return nil, err
	}

	// 3. AI Parse Preferences
	rules := []TeacherRule{}
	if combinedPrompt != "" {
		r, err := u.ScheduleAI.ParsePreferences(combinedPrompt, teachers)
		if err != nil {
			fmt.Printf("AI Parsing error: %v\n", err)
		} else {
			rules = r
		}
	}

	// 4. Schedule
	res, err := Can(teachers, classrooms, rules)
	if err != nil {
		return nil, err
	}

	// 4. Save
	for _, ss := range res.StudentScheduls {
		if _, err := u.ScadulStudentRepo.Create(&ss); err != nil {
			return nil, err
		}
	}
	for _, ts := range res.TeacherScheduls {
		if _, err := u.ScadulTeacherRepo.Create(&ts); err != nil {
			return nil, err
		}
	}

	return &res, nil
}
