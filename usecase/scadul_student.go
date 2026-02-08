package usecase

import (
	"scadulDataMono/domain/entities"
	"scadulDataMono/infra/gormDB/repo"
)

type ScadulStudentMg struct {
	ScadulStudentRepo *repo.ScadulStudentRepo
}

func NewScadulStudentMg(repo *repo.ScadulStudentRepo) *ScadulStudentMg {
	return &ScadulStudentMg{ScadulStudentRepo: repo}
}

func (u *ScadulStudentMg) Create(newS *entities.ScadulStudent) (uint, error) {
	return u.ScadulStudentRepo.Create(newS)
}

func (u *ScadulStudentMg) Update(id uint, updated *entities.ScadulStudent) (*entities.ScadulStudent, error) {
	return u.ScadulStudentRepo.Update(id, updated)
}

func (u *ScadulStudentMg) Delete(id uint) error {
	return u.ScadulStudentRepo.Delete(id)
}

func (u *ScadulStudentMg) Listing(search string, page, perPage int) ([]entities.ScadulStudent, int64, error) {
	return u.ScadulStudentRepo.Listing(search, page, perPage)
}

func (u *ScadulStudentMg) GetByID(id uint) (*entities.ScadulStudent, error) {
	return u.ScadulStudentRepo.GetByID(id)
}

func (u *ScadulStudentMg) AddSubjects(subjects []entities.SubjectInScadulStudent) error {
	return u.ScadulStudentRepo.AddSubjects(subjects)
}

func (u *ScadulStudentMg) RemoveSubjects(ids []uint) error {
	return u.ScadulStudentRepo.RemoveSubjects(ids)
}
