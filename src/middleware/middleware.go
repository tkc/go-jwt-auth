package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/tkc/go-jwt-auth/src/models"
	"github.com/tkc/go-jwt-auth/src/tokenizer"
	"net/http"
	"strings"
)

type middleware struct {
	tokenizer tokenizer.Tokenizer
}

func CreateMiddleware(tokenizer tokenizer.Tokenizer) Middleware {
	return &middleware{tokenizer}
}

type Middleware interface {
	IsLoggedIn() gin.HandlerFunc
	IsAdmin() gin.HandlerFunc
}

func extractToken(r *http.Request) *string {
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return &strArr[1]
	}
	return nil
}

func (m *middleware) IsLoggedIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := extractToken(c.Request)
		if tokenString == nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		claims, err := m.tokenizer.ParseToken(*tokenString)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Set(gin.AuthUserKey, models.User{
			ID:   cast.ToInt(claims["sub"]),
			Name: cast.ToString(claims["name"]),
		})
		c.Next()
	}
}

func (m *middleware) IsAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get(gin.AuthUserKey)
		if exists {
			c.JSON(http.StatusUnauthorized, "auth error")
			return
		}
		jwtUser := user.(*jwt.Token)
		claims := jwtUser.Claims.(jwt.MapClaims)
		if claims["role"] != "admin" {
			c.JSON(http.StatusUnauthorized, "auth error")
			return
		}
		c.Next()
	}
}
