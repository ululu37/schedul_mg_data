package usecase

import (
	"scadulDataMono/domain/entities"
	"scadulDataMono/infra/gormDB/repo"
)

type TeacherMg struct {
	teacherRepo *repo.TeacherRepo
	authRepo    *repo.AuthRepo
}

// NewTeacherMg creates a new TeacherMg instance
func NewTeacherMg(teacherRepo *repo.TeacherRepo, authRepo *repo.AuthRepo) *TeacherMg {
	return &TeacherMg{teacherRepo: teacherRepo, authRepo: authRepo}
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
	if len(existingSubjects) > 0 {
		//	return errors.New("one or more subjects already exist for this teacher")
		for is, _ := range subjects {
			for i, _ := range existingSubjects {
				if subjects[is].SubjectID == existingSubjects[i].SubjectID {
					existingSubjects = append(existingSubjects[:i], existingSubjects[i+1:]...)
					subjects = append(subjects[:is], subjects[is+1:]...)
					break
				}
			}
		}
		if len(subjects) == 0 {
			return nil //errors.New("no subjects to add")
		}
	}
	return u.teacherRepo.AddMySubject(teacherID, subjects)
}

func (u *TeacherMg) RemoveMySubject(ids []uint) error {
	return u.teacherRepo.RemoveMySubject(ids)
}

func (u *TeacherMg) GetMySubject(teacherID uint, minPreference int, page, perPage int) ([]entities.TeacherMySubject, int64, error) {
	return u.teacherRepo.GetMySubject(teacherID, minPreference, page, perPage)
}

func (u *TeacherMg) Listing(search string, page, perPage int) ([]entities.Teacher, int64, error) {
	return u.teacherRepo.Listing(search, page, perPage)
}
func (u *TeacherMg) Delete(id uint) error {
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
