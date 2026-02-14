package usecase

import (
	"scadulDataMono/domain/entities"
	"scadulDataMono/infra/gormDB/repo"
)

type TeacherMg struct {
	teacherRepo       *repo.TeacherRepo
	authRepo          *repo.AuthRepo
	scadulTeacherRepo *repo.ScadulTeacherRepo
}

// NewTeacherMg creates a new TeacherMg instance
func NewTeacherMg(teacherRepo *repo.TeacherRepo, authRepo *repo.AuthRepo, scadulTeacherRepo *repo.ScadulTeacherRepo) *TeacherMg {
	return &TeacherMg{
		teacherRepo:       teacherRepo,
		authRepo:          authRepo,
		scadulTeacherRepo: scadulTeacherRepo,
	}
}

func (u *TeacherMg) Create(name string, resume string, username string, password string, role int) (uint, error) {
	newAuth := &entities.Auth{
		Username:  username,
		Password:  password,
		HumanType: "t",
		Role:      role,
	}
	authID, err := u.authRepo.Create(newAuth)
	if err != nil {
		return 0, err
	}

	newT := &entities.Teacher{
		Name:   name,
		Resume: resume,
		AuthID: authID,
	}
	return u.teacherRepo.Create(newT)
}

func (u *TeacherMg) AddMySubject(teacherID uint, subjects []entities.TeacherMySubject) error {
	subjectIDs := make([]uint, len(subjects))
	for i, subject := range subjects {
		subjectIDs[i] = subject.SubjectID
	}

	// Check if any of the subjects already exist for this teacher
	existingSubjects, err := u.teacherRepo.GetMySubjectBySubjectIDs(teacherID, subjectIDs)
	if err != nil {
		return err
	}
	// Filter out already existing subjects
	filteredSubjects := []entities.TeacherMySubject{}
	existingMap := make(map[uint]bool)
	for _, es := range existingSubjects {
		existingMap[es.SubjectID] = true
	}

	for _, s := range subjects {
		if !existingMap[s.SubjectID] {
			filteredSubjects = append(filteredSubjects, s)
		}
	}
	subjects = filteredSubjects

	if len(subjects) == 0 {
		return nil
	}

	return u.teacherRepo.AddMySubject(teacherID, subjects)
}

func (u *TeacherMg) RemoveMySubject(ids []uint) error {
	return u.teacherRepo.RemoveMySubject(ids)
}

func (u *TeacherMg) DeleteAllMySubjects(teacherID uint) error {
	return u.teacherRepo.DeleteAllMySubjects(teacherID)
}

func (u *TeacherMg) GetMySubject(teacherID uint, minPreference int, search string, page, perPage int) ([]entities.TeacherMySubject, int64, error) {
	return u.teacherRepo.GetMySubject(teacherID, minPreference, search, page, perPage)
}

func (u *TeacherMg) Listing(search string, page, perPage int) ([]entities.Teacher, int64, error) {
	return u.teacherRepo.Listing(search, page, perPage)
}
func (u *TeacherMg) Delete(id uint) error {
	if err := u.teacherRepo.DeleteAllMySubjects(id); err != nil {
		return err
	}
	// Delete any schedule data associated with the teacher
	if u.scadulTeacherRepo != nil {
		if err := u.scadulTeacherRepo.DeleteByTeacherID(id); err != nil {
			return err
		}
	}
	return u.teacherRepo.Delete(id)
}

func (u *TeacherMg) DeleteByAuthID(authID uint) error {
	return u.teacherRepo.DeleteByAuthID(authID)
}
func (u *TeacherMg) EditPreference(updates []repo.PreferenceUpdate) error {
	return u.teacherRepo.UpdatePreference(updates)
}

func (u *TeacherMg) GetByID(teacherID uint) (*entities.Teacher, error) {
	return u.teacherRepo.GetByID(teacherID)
}

func (u *TeacherMg) Update(id uint, name, resume string) (*entities.Teacher, error) {
	updated := &entities.Teacher{
		Name:   name,
		Resume: resume,
	}
	return u.teacherRepo.Update(id, updated)
}
