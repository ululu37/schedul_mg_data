package repo

import (
	"scadulDataMono/domain/entities"
	"strconv"

	"gorm.io/gorm"
)

type TeacherRepo struct {
	DB *gorm.DB
}

func NewTeacherRepo(db *gorm.DB) *TeacherRepo {
	return &TeacherRepo{DB: db}
}
func (r *TeacherRepo) Create(newT *entities.Teacher) (uint, error) {
	if err := r.DB.Create(newT).Error; err != nil {
		return 0, err
	}
	return newT.ID, nil
}

func (r *TeacherRepo) Update(id uint, updated *entities.Teacher) (*entities.Teacher, error) {
	t := &entities.Teacher{}
	if err := r.DB.First(t, id).Error; err != nil {
		return nil, err
	}
	t.Name = updated.Name
	t.Resume = updated.Resume

	if err := r.DB.Save(t).Error; err != nil {
		return nil, err
	}

	// Reload with associations if needed, or just return t which now has updated fields.
	// However, if we want to be sure we have everything as a fresh get would return:
	if err := r.DB.Preload("Auth").First(t, id).Error; err != nil {
		return nil, err
	}

	return t, nil
}

func (r *TeacherRepo) Delete(id uint) error {
	res := r.DB.Delete(&entities.Teacher{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *TeacherRepo) DeleteByAuthID(authID uint) error {
	res := r.DB.Where("auth_id = ?", authID).Delete(&entities.Teacher{})
	if res.Error != nil {
		return res.Error
	}
	return nil
}

// Listing searches by teacher name, auth_id, or resume
func (r *TeacherRepo) Listing(search string, page, perPage int) ([]entities.Teacher, int64, error) {
	var list []entities.Teacher
	var count int64
	q := r.DB.Model(&entities.Teacher{})
	if search != "" {
		like := "%" + search + "%"
		q = q.Where("name LIKE ? OR resume LIKE ?", like, like)
		if authID, err := strconv.ParseUint(search, 10, 64); err == nil {
			q = q.Or("auth_id = ?", authID)
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

func (r *TeacherRepo) GetByIDs(ids []uint) ([]entities.Teacher, error) {
	var list []entities.Teacher
	if err := r.DB.Where("id IN ?", ids).Preload("Auth").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

type PreferenceUpdate struct {
	ID         uint
	Preference int
}

func (r *TeacherRepo) UpdatePreference(updates []PreferenceUpdate) error {
	if len(updates) == 0 {
		return nil
	}
	tx := r.DB.Begin()
	for _, u := range updates {
		res := tx.Model(&entities.TeacherMySubject{}).Where("id = ?", u.ID).Update("preference", u.Preference)
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

// GetMySubject returns TeacherMySubject for a teacher with min preference, paginated
func (r *TeacherRepo) GetMySubject(teacherID uint, minPreference int, page, perPage int) ([]entities.TeacherMySubject, int64, error) {
	var list []entities.TeacherMySubject
	var count int64
	q := r.DB.Model(&entities.TeacherMySubject{}).
		Where("teacher_id = ? AND preference >= ?", teacherID, minPreference).
		Preload("Subject")
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

// AddMySubject adds subjects to a teacher
func (r *TeacherRepo) AddMySubject(teacherID uint, subjects []entities.TeacherMySubject) error {
	if len(subjects) == 0 {
		return nil
	}
	for i := range subjects {
		subjects[i].TeacherID = teacherID
	}
	if err := r.DB.Create(&subjects).Error; err != nil {
		return err
	}
	return nil
}

// RemoveMySubject removes TeacherMySubject records by their IDs
func (r *TeacherRepo) RemoveMySubject(ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	res := r.DB.Delete(&entities.TeacherMySubject{}, ids)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// GetMySubjectBySubjectIDs returns TeacherMySubject records for given SubjectIDs
func (r *TeacherRepo) GetMySubjectBySubjectIDs(teacherID uint, subjectIDs []uint) ([]entities.TeacherMySubject, error) {
	var list []entities.TeacherMySubject
	if len(subjectIDs) == 0 {
		return list, nil
	}
	if err := r.DB.Where("teacher_id = ? AND subject_id IN ?", teacherID, subjectIDs).Preload("Subject").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// GetByID returns a Teacher entity by its ID (no MySubject)
func (r *TeacherRepo) GetByID(id uint) (*entities.Teacher, error) {
	var teacher entities.Teacher
	if err := r.DB.Where("id = ?", id).First(&teacher).Error; err != nil {
		return nil, err
	}
	return &teacher, nil
}

// NewTeacherRepo creates a new TeacherRepo instance
