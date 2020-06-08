package handlers

import (
	"net/http"
	"simple-api/middleware"

	"github.com/labstack/echo/v4"
)

type Error map[string]interface{}

var (
	InternalError = http.StatusInternalServerError
	BadRequest    = http.StatusBadRequest
)

type UserForm struct {
	ID       string `json:"id" uri:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (h *Handler) handleGetAllUsers() echo.HandlerFunc {
	return func(c echo.Context) error {
		users, err := h.DB.GetAllUsers()
		if err != nil {
			return c.JSON(InternalError, err)
		}
		return c.JSON(http.StatusOK, users)
	}
}
func (h *Handler) handleGetFollowersForUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		followers, err := h.DB.GetFollowersForUser(c.Get("userId").(string))
		if err != nil {
			return c.JSON(InternalError, err)
		}
		return c.JSON(http.StatusOK, followers)
	}
}

func (h *Handler) handleGetFollowingForUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		following, err := h.DB.GetFollowingForUser(c.Get("userId").(string))
		if err != nil {
			return c.JSON(InternalError, err)
		}
		return c.JSON(http.StatusOK, following)
	}

}

func (h *Handler) handleUnFollowUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		var user *UserForm
		if err := c.Bind(&user); err != nil {
			return c.JSON(BadRequest, err.Error())
		}
		if user.ID == "" {
			return c.JSON(BadRequest, "invalid input data")
		}
		if err := h.DB.UnfollowUser(c.Get("userId").(string), user.ID); err != nil {
			return c.JSON(InternalError, err.Error())
		}
		return nil
	}
}

func (h *Handler) handleFollowUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		var user *UserForm
		if err := c.Bind(&user); err != nil {
			return c.JSON(BadRequest, err.Error())
		}
		currentUserId := c.Get("user").(map[string]interface{})["id"].(string)
		if err := h.DB.FollowUser(currentUserId, user.ID); err != nil {
			h.E.Logger.Error(err)
			return c.JSON(InternalError, err.Error())
		}
		return nil
	}
}

func (h *Handler) handleCreateUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		var user *UserForm
		if err := c.Bind(&user); err != nil {
			return c.JSON(BadRequest, err.Error())
		}
		_, err := h.DB.CreateUser(user.Username, user.Password, user.Email)
		if err != nil {
			return c.JSON(InternalError, err.Error())
		}
		return c.NoContent(http.StatusCreated)
	}
}

func (h *Handler) handleDeleteUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		userId := c.Param("id")
		if err := h.DB.DeleteUser(userId); err != nil {
			return c.JSON(InternalError, err.Error())
		}
		return c.NoContent(http.StatusOK)
	}
}

func (h *Handler) handleLogIn() echo.HandlerFunc {
	return func(c echo.Context) error {
		var user UserForm
		if err := c.Bind(&user); err != nil {
			return c.JSON(http.StatusBadRequest, "failed to bind user")
		}
		authUser, err := h.DB.LogIn(user.Username, user.Password)
		if err != nil {
			return c.JSON(InternalError, err.Error())
		}

		t, err := middleware.GenerateToken(middleware.Claims{
			ID:       authUser.ID,
			Username: authUser.Username,
			Email:    authUser.Email,
		})
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, t)
	}
}

func (h *Handler) handleGetUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		userId := c.Get("user").(map[string]interface{})["id"].(string)
		user, err := h.DB.GetUser(userId)
		if err != nil {
			return c.JSON(InternalError, err.Error())
		}
		return c.JSON(200, user)
	}
}
