package tokenizer

import (
	"encoding/base64"
	"encoding/json"
	"github.com/bxcodec/faker/v3"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"github.com/tkc/go-jwt-auth/src/config"
	"github.com/tkc/go-jwt-auth/src/models"
	"strings"
	"testing"
	"time"
)

func normalToken() string {
	user := models.User{
		ID:   1,
		Role: models.RoleTypes.Admin,
		Name: faker.Name(),
	}

	tokenizer := CreateTokenizer()
	accessTokenClaims := jwt.MapClaims{
		"sub":   user.ID,
		"name":  user.Name,
		"admin": user.Role,
		"exp":   time.Now().Add(time.Minute * 15).Unix(), // 15min
		"iat":   time.Now().Unix(),
	}
	accessToken, _ := tokenizer.GenerateToken(accessTokenClaims)
	return *accessToken
}

func badTimeToken() string {
	token, _ := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub":  1,
			"name": faker.Name(),
			"exp":  time.Now().Add(time.Minute * -1).Unix(), // bad exp -1 min
		}).SignedString([]byte(config.SIGNINGKEY))
	return token
}

func badSignedToken() string {
	token, _ := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub":  1,
			"name": faker.Name(),
			"exp":  time.Now().Add(time.Minute * 15).Unix(),
		}).SignedString([]byte(faker.Sentence())) // bad signedS
	return token
}

func signingArgS384Token() string {
	token, _ := jwt.NewWithClaims(
		jwt.SigningMethodHS384, // bad arg S384
		jwt.MapClaims{
			"sub":  1,
			"name": faker.Name(),
			"exp":  time.Now().Add(time.Minute * 15).Unix(),
		}).SignedString([]byte(config.SIGNINGKEY))
	return token
}

func signingArgNoneToken() string {
	header := map[string]interface{}{
		"typ": "JWT",
		"alg": "none",
	}

	claims := map[string]interface{}{
		"sub": 1,
	}

	headerJsonValue, err := json.Marshal(header)
	if err != nil {
		return ""
	}

	claimsJsonValue, err := json.Marshal(claims)
	if err != nil {
		return ""
	}

	sstr := strings.Join([]string{
		strings.TrimRight(base64.URLEncoding.EncodeToString(headerJsonValue), "="),
		strings.TrimRight(base64.URLEncoding.EncodeToString(claimsJsonValue), "="),
	}, ".")

	token := strings.Join([]string{sstr, "dummy"}, ".")

	return token
}

func TestGenerateAccessToken(t *testing.T) {
	var tests = []struct {
		name  string
		input string
		isErr bool
	}{
		{"normalToken", normalToken(), false},
		{"badTimeToken", badTimeToken(), true},
		{"badSignedToken", badSignedToken(), true},
		{"signingArgS384Token", signingArgS384Token(), true},
		{"signingArgNoneToken", signingArgNoneToken(), true},
	}
	tokenizer := CreateTokenizer()
	for _, test := range tests {
		_, err := tokenizer.ParseToken(test.input)
		if test.isErr {
			assert.Error(t, err, test.name)
			continue
		}
		assert.Nil(t, err)
	}
}
