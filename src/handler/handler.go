package handler

import (
	"github.com/tkc/go-jwt-auth/src/models"
	"github.com/tkc/go-jwt-auth/src/tokenizer"
	"net/http"

	"github.com/gin-gonic/gin"
)

type handler struct {
	tokenizer tokenizer.Tokenizer
}

func CreateHandler(tokenizer tokenizer.Tokenizer) *handler {
	return &handler{tokenizer: tokenizer}
}

func (h *handler) Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	if username == "test" && password == "password" {
		tokens, err := h.tokenizer.GenerateTokenPair()
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
			return
		}
		c.JSON(http.StatusOK, tokens)
		return
	}
}

func (h *handler) Token(c *gin.Context) {
	type tokenReqBody struct {
		RefreshToken string `json:"refresh_token"`
	}

	tokenReq := tokenReqBody{}
	err := c.Bind(&tokenReq)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}

	claims, err := h.tokenizer.ParseToken(tokenReq.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}

	// Get the user record from database or
	// run through your business logic to verify if the user can log in
	if int(claims["sub"].(float64)) == 1 {
		newTokenPair, err := h.tokenizer.GenerateTokenPair()
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
			return
		}
		c.JSON(http.StatusOK, newTokenPair)
		return
	}
	c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
}

func (h *handler) Private(c *gin.Context) {
	user, exists := c.Get(gin.AuthUserKey)
	if exists {
		c.JSON(http.StatusUnauthorized, "auth error")
		return
	}
	c.String(http.StatusOK, "Welcome "+user.(models.User).Name+"!")
}
