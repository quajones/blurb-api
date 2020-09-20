package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type BlurbForm struct {
	Topic   int      `json:"topic`
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json: "tags"`
}

var (
	insertBlurbQuery = `insert into blurbs 
						(user_id, topic, title, content, created_date, last_modified_date) 
						values ($1, $2, $3, $4, $5, $6)`
	deleteBlurbQuery = `delete from blurbs where blurb_id = $1`
	createTagQuery   = `insert into tags (tag) values ($1)`
)

func (h *Handler) handleDeleteBlurb() echo.HandlerFunc {
	return func(c echo.Context) error {
		var blurbId string

	}
}
func (h *Handler) handleCreateBlurb() echo.HandlerFunc {
	return func(c echo.Context) error {
		var blurb *BlurbForm
		if err := c.Bind(&blurb); err != nil {
			return c.JSON(BadRequest, err.Error())
		}
		id, err := h.DB.CreateBlurb(c.Get("userId").(string), blurb.Topic, blurb.Title, blurb.Content)
		if err != nil {
			c.Logger().Error(err)
			return c.JSON(InternalError, err.Error())
		}
		for _, v := range blurb.Tags {
			if err := h.DB.CheckTag(id, v); err != nil {
				c.Logger().Error(err)
				return c.JSON(InternalError, err.Error())
			}
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
