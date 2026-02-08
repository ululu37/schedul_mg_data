package repo

import (
	"scadulDataMono/domain/entities"

	"gorm.io/gorm"
)

type TermRepo struct {
	DB *gorm.DB
}

func NewTermRepo(db *gorm.DB) *TermRepo {
	return &TermRepo{DB: db}
}

func (r *TermRepo) Create(term *entities.Term) (uint, error) {
	if err := r.DB.Create(term).Error; err != nil {
		return 0, err
	}
	return term.ID, nil
}

func (r *TermRepo) Listing(search string, page, perPage int) ([]entities.Term, int64, error) {
	var list []entities.Term
	var count int64
	q := r.DB.Model(&entities.Term{})

	if search != "" {
		q = q.Where("name LIKE ?", "%"+search+"%")
	}

	if err := q.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if page > 0 && perPage > 0 {
		offset := (page - 1) * perPage
		q = q.Limit(perPage).Offset(offset)
	}

	if err := q.Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, count, nil
}

func (r *TermRepo) GetByID(id uint) (*entities.Term, error) {
	term := &entities.Term{}
	if err := r.DB.First(term, id).Error; err != nil {
		return nil, err
	}
	return term, nil
}

func (r *TermRepo) Update(id uint, updated *entities.Term) (*entities.Term, error) {
	term := &entities.Term{}
	if err := r.DB.First(term, id).Error; err != nil {
		return nil, err
	}

	term.Name = updated.Name

	if err := r.DB.Save(term).Error; err != nil {
		return nil, err
	}
	return term, nil
}

func (r *TermRepo) Delete(id uint) error {
	res := r.DB.Delete(&entities.Term{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
