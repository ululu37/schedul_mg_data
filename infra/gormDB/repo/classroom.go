package repo

import (
	"scadulDataMono/domain/entities"
	"strconv"

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

	c.Name = updated.Name
	c.PreCurriculumID = updated.PreCurriculumID
	c.CurriculumID = updated.CurriculumID
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

// Listing searches by Name, PreCurriculum Name, Curriculum Name, or Year
func (r *ClassroomRepo) Listing(search string, page, perPage int) ([]entities.Classroom, int64, error) {
	var list []entities.Classroom
	var count int64

	// Start with base query
	q := r.DB.Model(&entities.Classroom{}).
		Joins("LEFT JOIN pre_curriculums ON pre_curriculums.id = classrooms.pre_curriculum_id").
		Joins("LEFT JOIN curriculums ON curriculums.id = classrooms.curriculum_id")

	if search != "" {
		like := "%" + search + "%"

		// Search in Classroom Name, PreCurriculum Name, Curriculum Name
		searchQuery := "classrooms.name LIKE ? OR pre_curriculums.name LIKE ? OR curriculums.name LIKE ?"
		args := []interface{}{like, like, like}

		// Also search in Year if search is numeric
		if year, err := strconv.Atoi(search); err == nil {
			searchQuery += " OR classrooms.year = ?"
			args = append(args, year)
		}

		q = q.Where(searchQuery, args...)
	}

	if err := q.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// Apply paging
	if page > 0 && perPage > 0 {
		offset := (page - 1) * perPage
		q = q.Limit(perPage).Offset(offset)
	}

	// Preload relations for the final list
	if err := q.Preload("PreCurriculum").Preload("Curriculum").Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, count, nil
}

func (r *ClassroomRepo) GetByID(id uint) (*entities.Classroom, error) {
	c := &entities.Classroom{}
	if err := r.DB.Preload("PreCurriculum").Preload("Curriculum").Preload("Student").First(c, id).Error; err != nil {
		return nil, err
	}
	return c, nil
}

func (r *ClassroomRepo) GetAllWithRelations() ([]entities.Classroom, error) {
	var list []entities.Classroom
	if err := r.DB.Preload("Curriculum.SubjectInCurriculum.SubjectInPreCurriculum.Subject").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
