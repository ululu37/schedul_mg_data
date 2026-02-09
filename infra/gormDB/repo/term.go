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

func (r *TermRepo) Listing() ([]entities.Term, error) {
	var list []entities.Term
	if err := r.DB.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
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
