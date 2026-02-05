package repo

import (
	"strconv"

	"scadulDataMono/domain/entities"

	"gorm.io/gorm"
)

type ClassroomRepo struct {
	DB *gorm.DB
}

func (r *ClassroomRepo) Create(newC *entities.Classroom) (uint, error) {
	if err := r.DB.Create(newC).Error; err != nil {
		return 0, err
	}
	return newC.ID, nil
}

func (r *ClassroomRepo) Update(id uint, updated *entities.Classroom) (*entities.Classroom, error) {
	c := &entities.Classroom{}
	if err := r.DB.First(c, id).Error; err != nil {
		return nil, err
	}

	c.TeacherID = updated.TeacherID
	c.Year = updated.Year

	if err := r.DB.Save(c).Error; err != nil {
		return nil, err
	}
	return c, nil
}

func (r *ClassroomRepo) Delete(id uint) error {
	res := r.DB.Delete(&entities.Classroom{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// Listing searches by classroom year or teacher name
func (r *ClassroomRepo) Listing(search string, page, perPage int) ([]entities.Classroom, int64, error) {
	var list []entities.Classroom
	var count int64
	q := r.DB.Model(&entities.Classroom{}).Preload("Teacher")
	if search != "" {
		like := "%" + search + "%"
		if id, err := strconv.ParseUint(search, 10, 64); err == nil {
			q = q.Where("year = ? OR teacher_id = ?", id, id)
		} else {
			q = q.Joins("LEFT JOIN teachers ON teachers.id = classrooms.teacher_id").Where("teachers.name LIKE ?", like)
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
