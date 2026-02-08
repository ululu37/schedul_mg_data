package usecase

import (
	"scadulDataMono/domain/entities"
	"scadulDataMono/infra/gormDB/repo"
)

type StudentMg struct {
	StudentRepo *repo.StudentRepo
	authRepo    *repo.AuthRepo
}

func NewStudentMg(repo *repo.StudentRepo, authRepo *repo.AuthRepo) *StudentMg {
	return &StudentMg{StudentRepo: repo, authRepo: authRepo}
}

func (u *StudentMg) Create(name string, year int, username string, password string, role int) (uint, error) {
	newAuth := &entities.Auth{
		Username:  username,
		Password:  password,
		HumanType: "s",
		Role:      role,
	}
	authID, err := u.authRepo.Create(newAuth)
	if err != nil {
		return 0, err
	}

	newT := &entities.Student{
		Name:   name,
		AuthID: authID,
		Year:   year,
	}
	return u.StudentRepo.Create(newT)
}
func (u *StudentMg) Listing(search string, page, perPage int) ([]entities.Student, int64, error) {
	return u.StudentRepo.Listing(search, page, perPage)
}

func (u *StudentMg) GetByID(id uint) (*entities.Student, error) {
	return u.StudentRepo.GetByID(id)
}

func (u *StudentMg) Update(id uint, updated *entities.Student) (*entities.Student, error) {
	return u.StudentRepo.Update(id, updated)
}

func (u *StudentMg) Delete(id uint) error {
	return u.StudentRepo.Delete(id)
}

func (u *StudentMg) DeleteByAuthID(authID uint) error {
	return u.StudentRepo.DeleteByAuthID(authID)
}
