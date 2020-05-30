package main

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

// Lesson contains info about the lesson of particular course
type Lesson struct {
	ID       int       `json:"id"`
	Title    string    `json:"title"`
	CourseID int       `json:"courseID"`
	Contents []Content `json:"contents"`
}

func createLesson(c echo.Context) error {
	tp := &Lesson{}
	if err := c.Bind(tp); err != nil {
		return err
	}
	lessons[tp.ID] = tp
	return c.JSON(http.StatusCreated, tp)
}

func getLessons(c echo.Context) error {
	var tp []Lesson = make([]Lesson, 0)

	for _, top := range lessons {
		tp = append(tp, *top)
	}
	return c.JSON(http.StatusOK, tp)
}

func getLessonByID(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	return c.JSON(http.StatusOK, lessons[id])
}

func updateLesson(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	top := &Lesson{}

	if err := c.Bind(&top); err != nil {
		return err
	}

	lessons[id] = top

	return c.JSON(http.StatusOK, top)
}

func deleteLesson(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	delete(lessons, id)

	return c.NoContent(http.StatusOK)

}
