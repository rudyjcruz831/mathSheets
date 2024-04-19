package handler

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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
	UserSevice       model.UserService
	TokenService     model.TokenService
	BaseURL          string
	TimeoutDurations time.Duration
	MaxBodyBytes     int64
}

// NewHandler creates a new instance of the Handler with the provided configuration.
func NewHandler(c *Config) {
	h := &Handler{
		UserService:  c.UserSevice,
		TokenService: c.TokenService,
	}

	// Enable CORS middleware with the provided configuration.
	c.R.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"POST"},
		AllowHeaders: []string{"Content-type"},
	}))

	g := c.R.Group(c.BaseURL) // Create a new route group under the base URL.

	// this is use to run test for CI and run code normaly on server
	// if gin.Mode() != gin.TestMode {

	// } else {

	// }

	// Define routes and their corresponding handler functions.
	g.GET("/", h.Home)               // Home route
	g.GET("/user/info")              // User info route (not implemented)
	g.POST("/user/signup", h.Signup) // User signup route
	g.POST("/user/signin", h.SignIn) // User signin route
}

func (h *Handler) Home(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]string{"Name": "Home"})
}

// func (h *Handler) user_info()
