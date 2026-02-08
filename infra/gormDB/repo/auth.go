package repo

import (
	"scadulDataMono/domain/entities"

	"gorm.io/gorm"
)

type AuthRepo struct {
	DB *gorm.DB
}

func NewAuthRepo(db *gorm.DB) *AuthRepo {
	return &AuthRepo{DB: db}
}

func (a *AuthRepo) Create(newAuth *entities.Auth) (uint, error) {
	if err := a.DB.Create(newAuth).Error; err != nil {
		return 0, err
	}
	return newAuth.ID, nil
}

func (a *AuthRepo) GetByUsername(username string) (*entities.Auth, error) {
	auth := &entities.Auth{}
	if err := a.DB.Where("username = ?", username).First(auth).Error; err != nil {
		return nil, err
	}
	return auth, nil
}

func (a *AuthRepo) Update(id uint, updatedAuth *entities.Auth) (*entities.Auth, error) {
	auth := &entities.Auth{}
	if err := a.DB.First(auth, id).Error; err != nil {
		return nil, err
	}

	// apply updates
	auth.Username = updatedAuth.Username
	auth.Password = updatedAuth.Password
	auth.HumanType = updatedAuth.HumanType
	auth.Role = updatedAuth.Role

	if err := a.DB.Save(auth).Error; err != nil {
		return nil, err
	}
	return auth, nil
}

func (a *AuthRepo) GetByID(id uint) (*entities.Auth, error) {
	auth := &entities.Auth{}
	if err := a.DB.First(auth, id).Error; err != nil {
		return nil, err
	}
	return auth, nil
}

func (a *AuthRepo) DeleteByID(id uint) error {
	res := a.DB.Delete(&entities.Auth{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (a *AuthRepo) Listing(search string, page, perPage int) ([]entities.Auth, int64, error) {
	var list []entities.Auth
	var count int64
	q := a.DB.Model(&entities.Auth{})

	if search != "" {
		q = q.Where("username LIKE ? OR human_type LIKE ?", "%"+search+"%", "%"+search+"%")
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
