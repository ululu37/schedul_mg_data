package repo

import (
	"scadulDataMono/domain/entities"

	"gorm.io/gorm"
)

type CurriculumRepo struct {
	DB *gorm.DB
}

func (r *CurriculumRepo) Create(newCur *entities.Curriculum) (uint, error) {
	if err := r.DB.Create(newCur).Error; err != nil {
		return 0, err
	}
	return newCur.ID, nil
}

func (r *CurriculumRepo) Update(id uint, updated *entities.Curriculum) (*entities.Curriculum, error) {
	cur := &entities.Curriculum{}
	if err := r.DB.First(cur, id).Error; err != nil {
		return nil, err
	}

	cur.Name = updated.Name
	if err := r.DB.Save(cur).Error; err != nil {
		return nil, err
	}
	return cur, nil
}

func (r *CurriculumRepo) Delete(id uint) error {
	res := r.DB.Delete(&entities.Curriculum{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *CurriculumRepo) Listing(name string, perPage, page int) ([]entities.Curriculum, int64, error) {
	var list []entities.Curriculum
	var count int64
	q := r.DB.Model(&entities.Curriculum{})
	if name != "" {
		q = q.Where("name LIKE ?", "%"+name+"%")
	}
	if err := q.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	offset := 0
	if page > 0 && perPage > 0 {
		offset = (page - 1) * perPage
	}
	if err := q.Preload("SubjectInCurriculum.SubjectInPreCurriculum.Subject").Limit(perPage).Offset(offset).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, count, nil
}

func (r *CurriculumRepo) GetByID(id uint, termID string) (*entities.Curriculum, error) {
	cur := &entities.Curriculum{}
	if termID == "" {
		if err := r.DB.Preload("SubjectInCurriculum.Term").Preload("SubjectInCurriculum.SubjectInPreCurriculum.Subject").First(cur, id).Error; err != nil {
			return nil, err
		}
		return cur, nil
	}
	// find term ID by name
	var t entities.Term
	if err := r.DB.Where("name = ?", termID).First(&t).Error; err != nil {
		return nil, err
	}
	if err := r.DB.Preload("SubjectInCurriculum", "term_id = ?", t.ID).Preload("SubjectInCurriculum.Term").Preload("SubjectInCurriculum.SubjectInPreCurriculum.Subject").First(cur, id).Error; err != nil {
		return nil, err
	}
	return cur, nil
}

type SubjectTermUpdate struct {
	ID     uint
	TermID uint
}

func (r *CurriculumRepo) AddSubject(subjects []entities.SubjectInCurriculum) error {
	if len(subjects) == 0 {
		return nil
	}
	if err := r.DB.Create(&subjects).Error; err != nil {
		return err
	}
	return nil
}

func (r *CurriculumRepo) RemoveSubject(ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	res := r.DB.Delete(&entities.SubjectInCurriculum{}, ids)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *CurriculumRepo) UpdateTerm(updates []SubjectTermUpdate) error {
	if len(updates) == 0 {
		return nil
	}
	tx := r.DB.Begin()
	for _, u := range updates {
		res := tx.Model(&entities.SubjectInCurriculum{}).Where("id = ?", u.ID).Update("term_id", u.TermID)
		if res.Error != nil {
			tx.Rollback()
			return res.Error
		}
		if res.RowsAffected == 0 {
			tx.Rollback()
			return gorm.ErrRecordNotFound
		}
	}
	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}
