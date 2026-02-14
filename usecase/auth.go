package usecase

import (
	"errors"
	"fmt"
	dto "scadulDataMono/domain/DTO"
	"scadulDataMono/domain/entities"
	"scadulDataMono/infra/gormDB/repo"
	jwthast "scadulDataMono/infra/jwt_hast"
)

type Auth struct {
	authRepo  *repo.AuthRepo
	studentMg *StudentMg
	teacherMg *TeacherMg
}

func NewAuth(authRepo *repo.AuthRepo, studentMg *StudentMg, teacherMg *TeacherMg) *Auth {
	return &Auth{
		authRepo:  authRepo,
		studentMg: studentMg,
		teacherMg: teacherMg,
	}
}

func (a *Auth) Login(username string, password string) (*dto.Passport, error) {
	fmt.Println("username", username)
	fmt.Println("password", password)
	auth, err := a.authRepo.GetByUsername(
		username,
	)
	if err != nil {
		return nil, err
	}
	if auth.Password != password {
		return nil, errors.New("Password invalid!")
	}

	payload := dto.PayLoad{
		ID:        auth.ID,
		HumanType: auth.HumanType,
		Role:      auth.Role,
	}
	token, err := jwthast.GenerateToken(payload)
	if err != nil {
		return nil, err
	}

	return &dto.Passport{
		Token:   token,
		Payload: payload,
	}, nil
}

func (a *Auth) Listing(search string, page, perPage int) ([]entities.Auth, int64, error) {
	return a.authRepo.Listing(search, page, perPage)
}

func (a *Auth) Update(id uint, updatedAuth *entities.Auth) (*entities.Auth, error) {
	return a.authRepo.Update(id, updatedAuth)
}

func (a *Auth) Delete(authID uint, humanType string) error {
	if humanType == "s" {
		if err := a.studentMg.DeleteByAuthID(authID); err != nil {
			return err
		}
	} else if humanType == "t" {
		if err := a.teacherMg.DeleteByAuthID(authID); err != nil {
			return err
		}
	}

	return a.authRepo.DeleteByID(authID)
}

func (a *Auth) GetProfile(authID uint, humanType string) (interface{}, error) {
	if humanType == "s" {
		return a.studentMg.GetByAuthID(authID)
	} else if humanType == "t" {
		return a.teacherMg.GetByAuthID(authID)
	}
	return nil, errors.New("admin has no human profile")
}
