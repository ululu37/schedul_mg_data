package usecase

import (
	"scadulDataMono/domain/entities"
	"scadulDataMono/infra/gormDB/repo"
)

type PreCurriculum struct {
	PreRepo   *repo.PreCuriculumRepo
	SubjectMg *SubjectMg
}

func (p *PreCurriculum) Create(name string) (uint, error) {
	pre := &entities.PreCurriculum{Name: name}
	return p.PreRepo.Create(pre)
}

func (p *PreCurriculum) Listing(search string, page, perPage int) ([]entities.PreCurriculum, int64, error) {
	return p.PreRepo.Listing(search, perPage, page)
}

func (p *PreCurriculum) Update(id uint, name string) (*entities.PreCurriculum, error) {
	updated := &entities.PreCurriculum{Name: name}
	return p.PreRepo.Update(id, updated)
}

func (p *PreCurriculum) Delete(id uint) error {
	return p.PreRepo.Delete(id)
}

func (p *PreCurriculum) CreateSubject(preCurriculumID uint, subjectName string, credit int) error {

	subjectID, err := p.SubjectMg.Create(subjectName)
	if err != nil {
		return err
	}

	return p.PreRepo.AddSubject([]entities.SubjectInPreCurriculum{{PreCurriculumID: preCurriculumID, SubjectID: subjectID, Credit: credit}})
}

func (p *PreCurriculum) RemoveSubject(SubjectInPreCurriculumID uint) error {
	subject, err := p.PreRepo.GetSubjectInPrecurriculumByID(SubjectInPreCurriculumID)
	if err != nil {
		return err
	}

	p.SubjectMg.Delete(subject.SubjectID)
	return p.PreRepo.RemoveSubject([]uint{SubjectInPreCurriculumID})
}
func (p *PreCurriculum) GetByID(id uint) (*entities.PreCurriculum, error) {
	return p.PreRepo.GetByID(id)
}
