package main

import (
	"github.com/gin-gonic/gin"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/tkc/go-jwt-auth/src/handler"
	"github.com/tkc/go-jwt-auth/src/middleware"
	"github.com/tkc/go-jwt-auth/src/tokenizer"

	"net/http"
)

var (
	router = gin.Default()
)

func main() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	tokenizer := tokenizer.CreateTokenizer()
	h := handler.CreateHandler(tokenizer)
	m := middleware.CreateMiddleware(tokenizer)

	router.POST("/login", h.Login)
	router.POST("/tokenizer", h.Token)

	adminRoute := router.Group("/admin")
	adminRoute.Use(
		m.IsLoggedIn(),
		m.IsAdmin())
	{
		adminRoute.GET("/", h.Private)
	}

	log.Fatal(router.Run(":5555"))

}
