package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type BlurbForm struct {
	Topic   int    `json:"topic`
	Title   string `json:"title"`
	Content string `json:"content"`
}

var (
	insertBlurbQuery = `insert into blurbs 
						(user_id, topic, title, content, created_date, last_modified_date) 
						values ($1, $2, $3, $4, $5, $6)`
	deleteBlurbQuery = `delete from blurbs where blurb_id = $1`
)

func (h *Handler) handleCreateBlurb() echo.HandlerFunc {
	fmt.Println("reached")
	return func(c echo.Context) error {
		var blurb *BlurbForm
		if err := c.Bind(&blurb); err != nil {
			return c.JSON(BadRequest, err.Error())
		}
		if err := h.DB.CreateBlurb(c.Get("userId").(string), blurb.Topic, blurb.Title, blurb.Content); err != nil {
			return c.JSON(InternalError, err.Error())
		}
		return c.NoContent(http.StatusCreated)
	}
}

func (h *Handler) handleGetBlurbsForUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		blurbs, err := h.DB.GetAllBlurbsForUser(c.Get("userId").(string))
		if err != nil {
			return c.JSON(InternalError, err.Error())
		}
		return c.JSON(200, blurbs)
	}
}

func (h *Handler) handleGetBlurbsForFollowing() echo.HandlerFunc {
	return func(c echo.Context) error {
		blurbs, err := h.DB.GetAllBlurbsForFollowing(c.Get("userId").(string))
		if err != nil {
			return c.JSON(InternalError, err.Error())
		}
		return c.JSON(200, blurbs)
	}
}
