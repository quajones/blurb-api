package handlers

import (
	_ "database/sql"

	"os"
	"simple-api/middleware"
	"simple-api/models"

	"github.com/labstack/echo/v4"
	m "github.com/labstack/echo/v4/middleware"
)

type Handler struct {
	DB models.DB
	E  *echo.Echo
}

func New() (*Handler, error) {
	e := echo.New()
	db, err := models.NewDB()
	if err != nil {
		return nil, err
	}
	return &Handler{db, e}, nil
}

func InitRoutes(h *Handler) {
	api := h.E.Group("/api")
	{
		api.POST("/signup", h.handleCreateUser())
		api.POST("/login", h.handleLogIn())
		api.GET("/users", h.handleGetAllUsers())
		api.Use(m.JWT([]byte(os.Getenv("BLURB_JWT"))), middleware.JWTCheck)
		api.GET("/user", h.handleGetUser())
		api.DELETE("/user/:id", h.handleDeleteUser())
		api.POST("/user/follow", h.handleFollowUser())
		api.POST("/user/unfollow", h.handleUnFollowUser())
		api.GET("/user/followers", h.handleGetFollowersForUser())
		api.GET("/user/following", h.handleGetFollowingForUser())
		api.POST("/blurb", h.handleCreateBlurb())
		api.GET("/blurb/:id", nil)
		api.GET("/user/blurbs", h.handleGetBlurbsForUser())
		api.GET("/u", h.handleGetBlurbsForFollowing())
		api.DELETE("/blurb/:id", nil)
		api.POST("/blurb/:id/quip", nil)
		api.DELETE("/blurb/:id/quip/:id", nil)
		api.POST("blurb/:id/clap", nil)
		api.POST("blurb/:id/unclap", nil)
	}
	blurbGroup := h.E.Group("/blurb")
	{
		blurbGroup.Use(m.JWT([]byte(os.Getenv("BLURB_JWT"))), middleware.JWTCheck)
		blurbGroup.POST("", h.handleCreateBlurb())
		blurbGroup.GET("/:id", nil)
		blurbGroup.GET("", h.handleGetBlurbsForUser())
		blurbGroup.GET("/u", h.handleGetBlurbsForFollowing())
		blurbGroup.DELETE("/:blurbId", nil)
		blurbGroup.POST("/clap/:blurbId", nil)
		blurbGroup.POST("/unclap/:blurbId", nil)
	}
}
