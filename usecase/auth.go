package usecase

import (
	"errors"
	dto "scadulDataMono/domain/DTO"
	"scadulDataMono/infra/gormDB/repo"
	jwthast "scadulDataMono/infra/jwt_hast"

	"github.com/golang-jwt/jwt/v5"
)

type Auth struct {
	authRepo repo.AuthRepo
}

func NewAuth(authRepo repo.AuthRepo) *Auth {
	return &Auth{authRepo: authRepo}
}

func (a *Auth) login(username string, password string) (*dto.Passport, error) {
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

type JwtCustomClaims struct {
	Payload dto.PayLoad
	jwt.RegisteredClaims
}
