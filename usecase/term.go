package usecase

import (
	"scadulDataMono/domain/entities"
	"scadulDataMono/infra/gormDB/repo"
)

type Term struct {
	Repo *repo.TermRepo
}

func NewTermUsecase(repo *repo.TermRepo) *Term {
	return &Term{Repo: repo}
}

func (u *Term) Create(name string) (uint, error) {
	term := &entities.Term{Name: name}
	return u.Repo.Create(term)
}

func (u *Term) Listing(search string, page, perPage int) ([]entities.Term, int64, error) {
	return u.Repo.Listing(search, page, perPage)
}

func (u *Term) GetByID(id uint) (*entities.Term, error) {
	return u.Repo.GetByID(id)
}

func (u *Term) Update(id uint, name string) (*entities.Term, error) {
	updated := &entities.Term{Name: name}
	return u.Repo.Update(id, updated)
}

func (u *Term) Delete(id uint) error {
	return u.Repo.Delete(id)
}
