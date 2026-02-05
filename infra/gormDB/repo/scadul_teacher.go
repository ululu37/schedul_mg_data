package repo

import (
	"scadulDataMono/domain/entities"

	"gorm.io/gorm"
)

type ScadulTeacherRepo struct {
	DB *gorm.DB
}

func (r *ScadulTeacherRepo) Create(newS *entities.ScadulTeacher) (uint, error) {
	if err := r.DB.Create(newS).Error; err != nil {
		return 0, err
	}
	return newS.ID, nil
}

func (r *ScadulTeacherRepo) Update(id uint, updated *entities.ScadulTeacher) (*entities.ScadulTeacher, error) {
	s := &entities.ScadulTeacher{}
	if err := r.DB.First(s, id).Error; err != nil {
		return nil, err
	}

	s.TeacherID = updated.TeacherID
	s.UseIn = updated.UseIn

	if err := r.DB.Save(s).Error; err != nil {
		return nil, err
	}
	return s, nil
}

func (r *ScadulTeacherRepo) Delete(id uint) error {
	res := r.DB.Delete(&entities.ScadulTeacher{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// Listing searches by useIn or teacher name
func (r *ScadulTeacherRepo) Listing(search string, page, perPage int) ([]entities.ScadulTeacher, int64, error) {
	var list []entities.ScadulTeacher
	var count int64
	q := r.DB.Model(&entities.ScadulTeacher{}).Preload("Teacher")
	if search != "" {
		like := "%" + search + "%"
		q = q.Joins("LEFT JOIN teachers ON teachers.id = scadul_teachers.teacher_id").Where("use_in LIKE ? OR teachers.name LIKE ?", like, like)
	}
	if err := q.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	offset := 0
	if page > 0 && perPage > 0 {
		offset = (page - 1) * perPage
		q = q.Limit(perPage).Offset(offset)
	}
	if err := q.Preload("SubjectInScadulTeacher.Teacher").Preload("SubjectInScadulTeacher.Subject").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, count, nil
}

func (r *ScadulTeacherRepo) GetByID(id uint) (*entities.ScadulTeacher, error) {
	s := &entities.ScadulTeacher{}
	if err := r.DB.Preload("SubjectInScadulTeacher.Teacher").Preload("SubjectInScadulTeacher.Subject").Preload("Teacher").First(s, id).Error; err != nil {
		return nil, err
	}
	return s, nil
}

func (r *ScadulTeacherRepo) AddSubjects(subjects []entities.SubjectInScadulTeacher) error {
	if len(subjects) == 0 {
		return nil
	}
	if err := r.DB.Create(&subjects).Error; err != nil {
		return err
	}
	return nil
}

func (r *ScadulTeacherRepo) RemoveSubjects(ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	res := r.DB.Delete(&entities.SubjectInScadulTeacher{}, ids)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
