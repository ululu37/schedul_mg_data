package jwthast

import (
	"errors"
	dto "scadulDataMono/domain/DTO"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtCustomClaims struct {
	Payload dto.PayLoad
	jwt.RegisteredClaims
}

func GenerateToken(payload dto.PayLoad) (string, error) {
	claims := JwtCustomClaims{
		Payload: payload,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("secret-key"))
}
func ParseToken(tokenString string) (*JwtCustomClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&JwtCustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte("secret-key"), nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JwtCustomClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
