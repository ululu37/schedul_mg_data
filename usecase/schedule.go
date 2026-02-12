package usecase

import (
	"scadulDataMono/infra/gormDB/repo"
)

type ScheduleUsecase struct {
	TeacherRepo       *repo.TeacherRepo
	ClassroomRepo     *repo.ClassroomRepo
	ScadulStudentRepo *repo.ScadulStudentRepo
	ScadulTeacherRepo *repo.ScadulTeacherRepo
}

func NewScheduleUsecase(
	teacherRepo *repo.TeacherRepo,
	classroomRepo *repo.ClassroomRepo,
	scadulStudentRepo *repo.ScadulStudentRepo,
	scadulTeacherRepo *repo.ScadulTeacherRepo,
) *ScheduleUsecase {
	return &ScheduleUsecase{
		TeacherRepo:       teacherRepo,
		ClassroomRepo:     classroomRepo,
		ScadulStudentRepo: scadulStudentRepo,
		ScadulTeacherRepo: scadulTeacherRepo,
	}
}

func (u *ScheduleUsecase) GenerateSchedule() (*SchedulRes, error) {
	// 0. Clear old schedules first
	if err := u.ScadulStudentRepo.DeleteAll(); err != nil {
		return nil, err
	}
	if err := u.ScadulTeacherRepo.DeleteAll(); err != nil {
		return nil, err
	}

	// 1. Load Data
	teachers, err := u.TeacherRepo.GetAllWithRelations()
	if err != nil {
		return nil, err
	}

	classrooms, err := u.ClassroomRepo.GetAllWithRelations()
	if err != nil {
		return nil, err
	}

	// 2. Schedule
	res, err := Can(teachers, classrooms)
	if err != nil {
		return nil, err
	}

	// 3. Save
	// Saving Student Schedules
	for _, ss := range res.StudentScheduls {
		// You might want to check if schedule exists for this classroom/term before creating
		if _, err := u.ScadulStudentRepo.Create(&ss); err != nil {
			return nil, err
		}
	}

	// Saving Teacher Schedules
	for _, ts := range res.TeacherScheduls {
		if _, err := u.ScadulTeacherRepo.Create(&ts); err != nil {
			return nil, err
		}
	}

	return &res, nil
}
