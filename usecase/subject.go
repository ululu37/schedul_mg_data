package usecase

import (
	"scadulDataMono/domain/entities"
	"scadulDataMono/infra/gormDB/repo" // Importing the repo package
)

type SubjectMg struct {
	SubjectRepo *repo.SubjectRepo // Changed to pointer to SubjectRepo
}

func (s *SubjectMg) Create(name []string) ([]uint, error) {
	subjects := make([]entities.Subject, 0, len(name))
	for _, n := range name {
		subjects = append(subjects, entities.Subject{Name: n})
	}
	ids, err := s.SubjectRepo.CreateMany(subjects)
	if err != nil {
		return nil, err
	}
	return ids, nil
}
func (s *SubjectMg) Update(id uint, updated *entities.Subject) (*entities.Subject, error) {
	return s.SubjectRepo.Update(id, updated) // New method for updating a subject
}

func (s *SubjectMg) Delete(id uint) error {
	return s.SubjectRepo.Delete(id) // New method for deleting a subject
}

func (s *SubjectMg) Listing(search string, page, perPage int) ([]entities.Subject, int64, error) {

	return s.SubjectRepo.Listing(search, page, perPage) // Updated to use SubjectRepo
}
