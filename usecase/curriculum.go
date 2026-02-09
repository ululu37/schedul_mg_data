package usecase

import (
	"scadulDataMono/domain/entities"
	"scadulDataMono/infra/gormDB/repo"
)

type CurriculumMg struct {
	CurriculumRepo *repo.CurriculumRepo
}

func NewCurriculumMg(repo *repo.CurriculumRepo) *CurriculumMg {
	return &CurriculumMg{CurriculumRepo: repo}
}

func (u *CurriculumMg) Create(name string, preCurriculumID uint) (uint, error) {
	newC := &entities.Curriculum{Name: name, PreCurriculumID: preCurriculumID}
	return u.CurriculumRepo.Create(newC)
}

func (u *CurriculumMg) Update(id uint, name string) (*entities.Curriculum, error) {
	updated := &entities.Curriculum{Name: name}
	return u.CurriculumRepo.Update(id, updated)
}

func (u *CurriculumMg) Delete(id uint) error {
	return u.CurriculumRepo.Delete(id)
}

func (u *CurriculumMg) Listing(search string, page, perPage int) ([]entities.Curriculum, int64, error) {
	return u.CurriculumRepo.Listing(search, perPage, page)
}

func (u *CurriculumMg) GetByID(id uint, termName string) (*entities.Curriculum, error) {
	return u.CurriculumRepo.GetByID(id, termName)
}

func (u *CurriculumMg) AddSubject(curriculumID uint, subjectPreIDs []uint) error {
	var subjects []entities.SubjectInCurriculum
	for _, id := range subjectPreIDs {
		subjects = append(subjects, entities.SubjectInCurriculum{
			CurriculumID:             curriculumID,
			SubjectInPreCurriculumID: id,
			TermID:                   nil,
		})
	}
	return u.CurriculumRepo.AddSubject(subjects)
}

func (u *CurriculumMg) EditSubjectTerm(updates []repo.SubjectTermUpdate) error {
	return u.CurriculumRepo.UpdateTerm(updates)
}

func (u *CurriculumMg) RemoveSubject(ids []uint) error {
	return u.CurriculumRepo.RemoveSubject(ids)
}
