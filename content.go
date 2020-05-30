package main

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

// Content is the core element. It has the actual course content.
type Content struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	LessonID    int    `json:"lessonID"`
	ContentType string `json:"contentType"`
	Data        string `json:"data"`
}

func createContent(c echo.Context) error {
	tp := &Content{}
	if err := c.Bind(tp); err != nil {
		return err
	}
	contents[tp.ID] = tp
	return c.JSON(http.StatusCreated, tp)
}

func getContents(c echo.Context) error {
	var tp []Content = make([]Content, 0)

	for _, top := range contents {
		tp = append(tp, *top)
	}
	return c.JSON(http.StatusOK, tp)
}

func getContentByID(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	return c.JSON(http.StatusOK, contents[id])
}

func updateContent(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	top := &Content{}

	if err := c.Bind(&top); err != nil {
		return err
	}

	contents[id] = top

	return c.JSON(http.StatusOK, top)
}

func deleteContent(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	delete(contents, id)

	return c.NoContent(http.StatusOK)

}
