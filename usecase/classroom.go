package usecase

import (
	"scadulDataMono/domain/entities"
	"scadulDataMono/infra/gormDB/repo"
)

type ClassroomMg struct {
	ClassroomRepo *repo.ClassroomRepo
}

func NewClassroomMg(repo *repo.ClassroomRepo) *ClassroomMg {
	return &ClassroomMg{ClassroomRepo: repo}
}

func (u *ClassroomMg) Create(newC *entities.Classroom) (uint, error) {
	return u.ClassroomRepo.Create(newC)
}

func (u *ClassroomMg) Update(id uint, updated *entities.Classroom) (*entities.Classroom, error) {
	return u.ClassroomRepo.Update(id, updated)
}

func (u *ClassroomMg) Delete(id uint) error {
	return u.ClassroomRepo.Delete(id)
}

func (u *ClassroomMg) Listing(search string, page, perPage int) ([]entities.Classroom, int64, error) {
	return u.ClassroomRepo.Listing(search, page, perPage)
}

func (u *ClassroomMg) GetByID(id uint) (*entities.Classroom, error) {
	return u.ClassroomRepo.GetByID(id)
}