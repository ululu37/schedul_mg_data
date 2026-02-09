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

func (u *Term) Listing() ([]entities.Term, error) {
	return u.Repo.Listing()
}

func (u *Term) Update(id uint, name string) (*entities.Term, error) {
	updated := &entities.Term{Name: name}
	return u.Repo.Update(id, updated)
}

func (u *Term) Delete(id uint) error {
	return u.Repo.Delete(id)
}
