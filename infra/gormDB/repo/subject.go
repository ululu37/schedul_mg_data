package repo

import (
	"scadulDataMono/domain/entities"

	"gorm.io/gorm"
)

type SubjectRepo struct {
	DB *gorm.DB
}

// NewSubjectRepo creates a new SubjectRepo instance
func NewSubjectRepo(db *gorm.DB) *SubjectRepo {
	return &SubjectRepo{DB: db}
}
func (r *SubjectRepo) Create(newSub *entities.Subject) (uint, error) {
	if err := r.DB.Create(newSub).Error; err != nil {
		return 0, err
	}
	return newSub.ID, nil
}

func (r *SubjectRepo) Update(id uint, updated *entities.Subject) (*entities.Subject, error) {
	sub := &entities.Subject{}
	if err := r.DB.First(sub, id).Error; err != nil {
		return nil, err
	}

	sub.Name = updated.Name
	if err := r.DB.Save(sub).Error; err != nil {
		return nil, err
	}
	return sub, nil
}

func (r *SubjectRepo) Delete(id uint) error {
	res := r.DB.Delete(&entities.Subject{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// Listing returns subjects filtered by search (name, preCurrilulumName, CuriiculumName) with pagination
func (r *SubjectRepo) Listing(search string, page, perPage int) ([]entities.Subject, int64, error) {
	var list []entities.Subject
	var count int64
	q := r.DB.Model(&entities.Subject{})
	if search != "" {
		like := "%" + search + "%"
		q = q.Where("name LIKE ? OR preCurrilulumName LIKE ? OR CuriiculumName LIKE ?", like, like, like)
	}
	if err := q.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	offset := 0
	if page > 0 && perPage > 0 {
		offset = (page - 1) * perPage
		q = q.Limit(perPage).Offset(offset)
	}
	if err := q.Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, count, nil
}
