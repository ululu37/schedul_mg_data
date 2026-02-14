package usecase

import (
	"scadulDataMono/domain/entities"
	"scadulDataMono/infra/gormDB/repo" // Importing the repo package
)

type SubjectMg struct {
	SubjectRepo *repo.SubjectRepo // Changed to pointer to SubjectRepo
}

func (s *SubjectMg) Create(names []string) ([]uint, error) {
	ids := make([]uint, len(names))
	for i, name := range names {
		// Check if exists
		existing, err := s.SubjectRepo.GetByName(name)
		if err == nil {
			ids[i] = existing.ID
			continue
		}

		// Create new
		newIDs, err := s.SubjectRepo.CreateMany([]entities.Subject{{Name: name}})
		if err != nil {
			return nil, err
		}
		ids[i] = newIDs[0]
	}
	return ids, nil
}
func (s *SubjectMg) Update(id uint, updated *entities.Subject) (*entities.Subject, error) {
	return s.SubjectRepo.Update(id, updated) // New method for updating a subject
}

func (s *SubjectMg) Delete(id uint) error {
	has, err := s.SubjectRepo.HasReferences(id)
	if err != nil {
		return err
	}
	// If still referenced by some curriculum or schedule, don't delete the core subject
	if has {
		return nil
	}
	return s.SubjectRepo.Delete(id)
}

func (s *SubjectMg) Listing(search string, page, perPage int) ([]entities.Subject, int64, error) {

	return s.SubjectRepo.Listing(search, page, perPage) // Updated to use SubjectRepo
}

func (s *SubjectMg) GetByID(id uint) (*entities.Subject, error) {
	return s.SubjectRepo.GetByID(id)
}
