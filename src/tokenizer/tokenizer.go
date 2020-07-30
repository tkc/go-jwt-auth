package tokenizer

import (
	"github.com/bxcodec/faker/v3"
	"github.com/tkc/go-jwt-auth/src/config"
	"github.com/tkc/go-jwt-auth/src/models"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

const (
	ALGORITHM = "HS256"
)

type tokenizer struct{}

type Tokenizer interface {
	GenerateTokenPair() (map[string]string, error)
	ParseToken(tokenString string) (jwt.MapClaims, error)
	GenerateToken(mapClaims jwt.MapClaims) (*string, error)
}

func CreateTokenizer() Tokenizer {
	return &tokenizer{}
}

func (h *tokenizer) GenerateTokenPair() (map[string]string, error) {
	user := models.User{
		ID:   1,
		Role: models.RoleTypes.Admin,
		Name: faker.Name(),
	}

	accessTokenClaims := baseClimes()
	accessTokenClaims["sib"] = user.ID
	accessTokenClaims["name"] = user.Name
	accessTokenClaims["role"] = user.Role
	accessTokenClaims["exp"] = time.Now().Add(time.Minute * 15).Unix() // 15min
	accessToken, err := h.GenerateToken(accessTokenClaims)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	refreshTokenClaims := baseClimes()
	refreshTokenClaims["sib"] = user.ID
	refreshTokenClaims["exp"] = time.Now().Add(time.Hour * 24).Unix() // 24h
	refreshToken, err := h.GenerateToken(refreshTokenClaims)

	if err != nil {
		return nil, errors.WithStack(err)
	}
	return map[string]string{"access_token": *accessToken, "refresh_token": *refreshToken}, nil
}

func (h *tokenizer) GenerateToken(claims jwt.MapClaims) (*string, error) {
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(config.SIGNINGKEY))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &token, nil
}

func (h *tokenizer) ParseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.WithStack(errors.New("invalid arg"))
			}
			return []byte(config.SIGNINGKEY), nil
		})

	if err != nil {
		return nil, errors.WithStack(err)
	}

	// Check Alg
	if token.Method.Alg() != ALGORITHM {
		return nil, errors.WithStack(errors.New("invalid arg"))
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok && !token.Valid {
		return nil, errors.WithStack(errors.New("Valid fail"))
	}
	return claims, nil
}

func baseClimes() jwt.MapClaims {
	claims := jwt.MapClaims{
		"kid":   "", // TODO
		"nonce": "", // TODO
		"iat":   time.Now().Unix(),
	}
	return claims
}
