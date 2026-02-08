package repo

import (
	"strconv"
	"strings"

	"scadulDataMono/domain/entities"

	"gorm.io/gorm"
)

type StudentRepo struct {
	DB *gorm.DB
}

func (r *StudentRepo) Create(newS *entities.Student) (uint, error) {
	if err := r.DB.Create(newS).Error; err != nil {
		return 0, err
	}
	return newS.ID, nil
}

func (r *StudentRepo) Update(id uint, updated *entities.Student) (*entities.Student, error) {
	s := &entities.Student{}
	if err := r.DB.First(s, id).Error; err != nil {
		return nil, err
	}

	s.Name = updated.Name
	//s.AuthID = updated.AuthID
	s.CurriculumID = updated.CurriculumID
	s.Year = updated.Year
	s.ClassroomID = updated.ClassroomID

	if err := r.DB.Save(s).Error; err != nil {
		return nil, err
	}

	// Refresh data with preloads
	if err := r.DB.Preload("Auth").Preload("Curriculum").Preload("Classroom").First(s, id).Error; err != nil {
		return nil, err
	}

	return s, nil
}

func (r *StudentRepo) Delete(id uint) error {
	res := r.DB.Delete(&entities.Student{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *StudentRepo) DeleteByAuthID(authID uint) error {
	res := r.DB.Where("auth_id = ?", authID).Delete(&entities.Student{})
	if res.Error != nil {
		return res.Error
	}
	return nil
}

// Listing searches by student name, student id, curriculum name, pre-curiculum name, or classroom teacher name
func (r *StudentRepo) Listing(search string, page, perPage int) ([]entities.Student, int64, error) {
	var list []entities.Student
	var count int64
	q := r.DB.Model(&entities.Student{}).
		Preload("Auth").
		Preload("Curriculum").
		Preload("Classroom").
		Joins("LEFT JOIN curriculums ON curriculums.id = students.curriculum_id").
		Joins("LEFT JOIN subject_in_curriculums ON subject_in_curriculums.curriculum_id = curriculums.id").
		Joins("LEFT JOIN subject_in_pre_curriculums ON subject_in_pre_curriculums.id = subject_in_curriculums.subject_in_pre_curriculum_id").
		Joins("LEFT JOIN pre_curriculums ON pre_curriculums.id = subject_in_pre_curriculums.pre_curriculum_id").
		Joins("LEFT JOIN classrooms ON classrooms.id = students.classroom_id").
		Joins("LEFT JOIN teachers ON teachers.id = classrooms.teacher_id")

	if strings.TrimSpace(search) != "" {
		if id, err := strconv.ParseUint(search, 10, 64); err == nil {
			q = q.Where("students.id = ? OR students.name LIKE ? OR curriculums.name LIKE ? OR pre_curriculums.name LIKE ? OR teachers.name LIKE ?", id, "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")
		} else {
			q = q.Where("students.name LIKE ? OR curriculums.name LIKE ? OR pre_curriculums.name LIKE ? OR teachers.name LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")
		}
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

func (r *StudentRepo) GetByClassRoom(classroomID uint) ([]entities.Student, error) {
	var list []entities.Student
	if err := r.DB.Where("classroom_id = ?", classroomID).Preload("Auth").Preload("Curriculum").Preload("Classroom").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *StudentRepo) GetByID(id uint) (*entities.Student, error) {
	s := &entities.Student{}
	if err := r.DB.Preload("Auth").Preload("Curriculum").Preload("Classroom").First(s, id).Error; err != nil {
		return nil, err
	}
	return s, nil
}
