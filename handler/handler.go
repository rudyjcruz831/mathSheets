package handler

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rudyjcruz831/mathSheets/model"
)

type Handler struct {
	UserService  model.UserService
	TokenService model.TokenService
	MaxBodyBytes int64
}

type Config struct {
	R                *gin.Engine
	UserSevice       model.UserService
	TokenService     model.TokenService
	BaseURL          string
	TimeoutDurations time.Duration
	MaxBodyBytes     int64
}

func NewHandler(c *Config) {
	h := &Handler{
		UserService:  c.UserSevice,
		TokenService: c.TokenService,
	}

	c.R.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"POST"},
		AllowHeaders: []string{"Content-type"},
	}))

	g := c.R.Group(c.BaseURL)

	// this is use to run test for CI/CD and run code normaly on server
	// if gin.Mode() != gin.TestMode {

	// } else {

	// }

	g.GET("/", h.Home)
	g.GET("/user/info")
	g.POST("/user/signup", h.Signup)
	g.POST("/user/signin", h.SignIn)
}

func (h *Handler) Home(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]string{"Name": "Home"})
}

// func (h *Handler) user_info()
