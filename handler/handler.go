package handler

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rudyjcruz831/mathSheets/handler/middleware"
	"github.com/rudyjcruz831/mathSheets/model"
)

// Handler is a struct that handles HTTP requests and responses.
type Handler struct {
	UserService  model.UserService
	TokenService model.TokenService
	MaxBodyBytes int64
}

// Config holds the configuration parameters for the Handler.
type Config struct {
	R                *gin.Engine
	UserService      model.UserService
	TokenService     model.TokenService
	BaseURL          string
	TimeoutDurations time.Duration
	MaxBodyBytes     int64
}

// NewHandler creates a new instance of the Handler with the provided configuration.
func NewHandler(c *Config) {
	h := &Handler{
		UserService:  c.UserService,
		TokenService: c.TokenService,
	}

	// Enable CORS middleware with the provided configuration.
	c.R.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"POST, GET"},
		AllowHeaders: []string{"Content-type, Authorization"},
	}))

	g := c.R.Group(c.BaseURL) // Create a new route group under the base URL.

	// this is use to run test for CI and run code normaly on local machine
	// if gin.Mode() != gin.TestMode {

	// } else {

	// }

	// Define routes and their corresponding handler functions.
	g.GET("/", h.Home) // Home route
	g.GET("/user/signout", middleware.AuthUser(h.TokenService, h.UserService), h.SignOut)
	g.GET("/user/info", middleware.AuthUser(h.TokenService, h.UserService), h.UserInfo)
	g.POST("/user/signup", h.Signup) // User signup route
	g.POST("/user/signin", h.SignIn) // User signin route
	g.POST("/user/tokens", h.Tokens)
}

func (h *Handler) Home(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]string{"Name": "Home"})
}

// func (h *Handler) user_info()
