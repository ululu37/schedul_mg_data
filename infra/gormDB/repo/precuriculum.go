package repo

import (
	"fmt"
	"scadulDataMono/domain/entities"

	"gorm.io/gorm"
)

type PreCuriculumRepo struct {
	DB *gorm.DB
}

func (r *PreCuriculumRepo) Create(newPre *entities.PreCurriculum) (uint, error) {
	if err := r.DB.Create(newPre).Error; err != nil {
		return 0, err
	}
	return newPre.ID, nil
}

func (r *PreCuriculumRepo) Update(id uint, updated *entities.PreCurriculum) (*entities.PreCurriculum, error) {
	p := &entities.PreCurriculum{}
	if err := r.DB.First(p, id).Error; err != nil {
		return nil, err
	}
	p.Name = updated.Name
	if err := r.DB.Save(p).Error; err != nil {
		return nil, err
	}
	return p, nil
}

func (r *PreCuriculumRepo) Delete(id uint) error {
	res := r.DB.Delete(&entities.PreCurriculum{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *PreCuriculumRepo) Listing(name string, perPage, page int) ([]entities.PreCurriculum, int64, error) {
	var list []entities.PreCurriculum
	var count int64
	q := r.DB.Model(&entities.PreCurriculum{})
	// Search by curriculum name
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
	if err := q.Limit(perPage).Offset(offset).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	fmt.Println("List:", list)
	return list, count, nil
}

func (r *PreCuriculumRepo) GetByID(id uint) (*entities.PreCurriculum, error) {
	p := &entities.PreCurriculum{}
	if err := r.DB.Preload("SubjectInPreCurriculum.Subject").First(p, id).Error; err != nil {
		return nil, err
	}
	fmt.Println("PreCurriculum:", p)
	return p, nil
}

func (r *PreCuriculumRepo) AddSubject(subjects []entities.SubjectInPreCurriculum) error {
	if len(subjects) == 0 {
		return nil
	}
	if err := r.DB.Create(&subjects).Error; err != nil {
		return err
	}
	return nil
}

func (r *PreCuriculumRepo) RemoveSubject(ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	res := r.DB.Delete(&entities.SubjectInPreCurriculum{}, ids)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// GetSubjectInPrecurriculumByID fetches a single SubjectInPreCurriculum by its ID
func (r *PreCuriculumRepo) GetSubjectInPrecurriculumByID(id uint) (*entities.SubjectInPreCurriculum, error) {
	var subject entities.SubjectInPreCurriculum
	if err := r.DB.Preload("Subject").First(&subject, id).Error; err != nil {
		return nil, err
	}
	return &subject, nil
}
