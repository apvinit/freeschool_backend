package main

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

// Course holds the info about a particular course
type Course struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	CategoryID  int    `json:"categoryID,omitempty"`
	IconURL     string `json:"iconURL,omitempty"`
}

func createCourse(c echo.Context) error {

	cou := &Course{}
	if err := c.Bind(cou); err != nil {
		return err
	}
	courses[cou.ID] = cou
	return c.JSON(http.StatusCreated, cou)
}

func getCourses(c echo.Context) error {
	var cou []Course = make([]Course, 0)

	for _, ct := range courses {
		cou = append(cou, *ct)
	}
	return c.JSON(http.StatusOK, cou)
}

func getCourseByID(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	return c.JSON(http.StatusOK, courses[id])
}

func updateCourse(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	cou := &Course{}

	if err := c.Bind(&cou); err != nil {
		return err
	}

	courses[id] = cou

	return c.JSON(http.StatusOK, cou)
}

func deleteCourse(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	delete(courses, id)

	return c.NoContent(http.StatusOK)
}
