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

func (p *PreCurriculum) CreateSubject(preCurriculumID uint, newSubjectInCurriculum []entities.SubjectInPreCurriculum) error {
	// 1. Get Global IDs (Create if not exist, or return existing ID)
	subjectNames := make([]string, 0, len(newSubjectInCurriculum))
	for _, s := range newSubjectInCurriculum {
		subjectNames = append(subjectNames, s.Subject.Name)
	}

	ids, err := p.SubjectMg.Create(subjectNames)
	if err != nil {
		return err
	}

	// 2. Fetch existing subjects in this curriculum to prevent local duplicates
	existingPre, err := p.PreRepo.GetByID(preCurriculumID)
	if err != nil {
		return err
	}
	existingMap := make(map[uint]bool)
	for _, s := range existingPre.SubjectInPreCurriculum {
		existingMap[s.SubjectID] = true
	}

	// 3. Prepare data and filter out local duplicates
	toAdd := []entities.SubjectInPreCurriculum{}
	for i := range newSubjectInCurriculum {
		newSubjectInCurriculum[i].SubjectID = ids[i]
		// Clear local Subject struct to avoid GORM conflicts
		newSubjectInCurriculum[i].Subject = entities.Subject{}

		if !existingMap[ids[i]] {
			toAdd = append(toAdd, newSubjectInCurriculum[i])
			// Mark as added to prevent duplicates within the same batch request
			existingMap[ids[i]] = true
		}
	}

	// 4. Perform batch add for new subjects only
	if len(toAdd) == 0 {
		return nil
	}

	return p.PreRepo.AddSubject(toAdd)
}

func (p *PreCurriculum) RemoveSubject(SubjectInPreCurriculumID uint) error {
	subject, err := p.PreRepo.GetSubjectInPrecurriculumByID(SubjectInPreCurriculumID)
	if err != nil {
		return err
	}

	// 1. Remove the link from PreCurriculum first
	if err := p.PreRepo.RemoveSubject([]uint{SubjectInPreCurriculumID}); err != nil {
		return err
	}

	// 2. Then try to clean up the actual Subject record (Delete if no other references)
	return p.SubjectMg.Delete(subject.SubjectID)
}
func (p *PreCurriculum) GetByID(id uint) (*entities.PreCurriculum, error) {
	return p.PreRepo.GetByID(id)
}
func (p *PreCurriculum) UpdateSubject(SubjectInPreCurriculumID uint, subjectName string, credit int) error {
	return p.PreRepo.UpdateSubjectInPreCurriculum(SubjectInPreCurriculumID, subjectName, credit)
}
