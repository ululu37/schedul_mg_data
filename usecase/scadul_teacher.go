package usecase

import (
	"scadulDataMono/domain/entities"
	"scadulDataMono/infra/gormDB/repo"
)

type ScadulTeacherMg struct {
	ScadulTeacherRepo *repo.ScadulTeacherRepo
}

func NewScadulTeacherMg(repo *repo.ScadulTeacherRepo) *ScadulTeacherMg {
	return &ScadulTeacherMg{ScadulTeacherRepo: repo}
}

func (u *ScadulTeacherMg) Create(newS *entities.ScadulTeacher) (uint, error) {
	return u.ScadulTeacherRepo.Create(newS)
}

func (u *ScadulTeacherMg) Update(id uint, updated *entities.ScadulTeacher) (*entities.ScadulTeacher, error) {
	return u.ScadulTeacherRepo.Update(id, updated)
}

func (u *ScadulTeacherMg) Delete(id uint) error {
	return u.ScadulTeacherRepo.Delete(id)
}

func (u *ScadulTeacherMg) Listing(search string, page, perPage int) ([]entities.ScadulTeacher, int64, error) {
	return u.ScadulTeacherRepo.Listing(search, page, perPage)
}

func (u *ScadulTeacherMg) GetByID(id uint) (*entities.ScadulTeacher, error) {
	return u.ScadulTeacherRepo.GetByID(id)
}

func (u *ScadulTeacherMg) AddSubjects(subjects []entities.SubjectInScadulTeacher) error {
	return u.ScadulTeacherRepo.AddSubjects(subjects)
}

func (u *ScadulTeacherMg) RemoveSubjects(ids []uint) error {
	return u.ScadulTeacherRepo.RemoveSubjects(ids)
}
