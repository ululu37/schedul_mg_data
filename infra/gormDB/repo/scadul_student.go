package repo

import (
	"scadulDataMono/domain/entities"

	"gorm.io/gorm"
)

type ScadulStudentRepo struct {
	DB *gorm.DB
}

func (r *ScadulStudentRepo) Create(newS *entities.ScadulStudent) (uint, error) {
	if err := r.DB.Create(newS).Error; err != nil {
		return 0, err
	}
	return newS.ID, nil
}

func (r *ScadulStudentRepo) Update(id uint, updated *entities.ScadulStudent) (*entities.ScadulStudent, error) {
	s := &entities.ScadulStudent{}
	if err := r.DB.First(s, id).Error; err != nil {
		return nil, err
	}

	s.ClassroomID = updated.ClassroomID
	s.UseIn = updated.UseIn

	if err := r.DB.Save(s).Error; err != nil {
		return nil, err
	}
	return s, nil
}

func (r *ScadulStudentRepo) Delete(id uint) error {
	res := r.DB.Delete(&entities.ScadulStudent{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// Listing searches by useIn or classroom id
func (r *ScadulStudentRepo) Listing(search string, page, perPage int) ([]entities.ScadulStudent, int64, error) {
	var list []entities.ScadulStudent
	var count int64
	q := r.DB.Model(&entities.ScadulStudent{}).Preload("Classroom")
	if search != "" {
		q = q.Where("use_in LIKE ? OR classroom_id = ?", "%"+search+"%", search)
	}
	if err := q.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	offset := 0
	if page > 0 && perPage > 0 {
		offset = (page - 1) * perPage
		q = q.Limit(perPage).Offset(offset)
	}
	if err := q.Preload("SubjectInScadulStudent.Teacher").Preload("SubjectInScadulStudent.Subject").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, count, nil
}

func (r *ScadulStudentRepo) GetByID(id uint) (*entities.ScadulStudent, error) {
	s := &entities.ScadulStudent{}
	if err := r.DB.Preload("SubjectInScadulStudent.Teacher").Preload("SubjectInScadulStudent.Subject").Preload("Classroom").First(s, id).Error; err != nil {
		return nil, err
	}
	return s, nil
}

func (r *ScadulStudentRepo) AddSubjects(subjects []entities.SubjectInScadulStudent) error {
	if len(subjects) == 0 {
		return nil
	}
	if err := r.DB.Create(&subjects).Error; err != nil {
		return err
	}
	return nil
}

func (r *ScadulStudentRepo) RemoveSubjects(ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	res := r.DB.Delete(&entities.SubjectInScadulStudent{}, ids)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *ScadulStudentRepo) DeleteAll() error {
	// Delete all subject assignments first
	if err := r.DB.Exec("DELETE FROM subject_in_scadul_students").Error; err != nil {
		return err
	}
	// Then delete all schedules
	if err := r.DB.Exec("DELETE FROM scadul_students").Error; err != nil {
		return err
	}
	return nil
}
