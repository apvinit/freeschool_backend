package main

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

// Course holds the info about a particular course
type Course struct {
	ID          int    `json:"id,omitempty"`
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

	insertCourseSQL := "INSERT INTO course(title, category_id) VALUES(?,?)"

	stmt, err := db.Prepare(insertCourseSQL)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(cou.Title, cou.CategoryID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, cou)
}

func getCourses(c echo.Context) error {
	var cou []Course = make([]Course, 0)

	row, err := db.Query("SELECT id, title, category_id FROM course")
	if err != nil {
		return err
	}
	defer row.Close()

	for row.Next() {
		co := Course{}
		row.Scan(&co.ID, &co.Title, &co.CategoryID)
		cou = append(cou, co)
	}

	return c.JSON(http.StatusOK, cou)
}

func getCourseByID(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	row := db.QueryRow("SELECT id, title, category_id FROM course WHERE id=?", id)
	co := Course{}
	row.Scan(&co.ID, &co.Title, &co.CategoryID)

	return c.JSON(http.StatusOK, co)
}

func updateCourse(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	cou := &Course{}

	if err := c.Bind(&cou); err != nil {
		return err
	}

	stmt, err := db.Prepare("UPDATE course SET title=? WHERE id=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(cou.Title, id)

	return c.JSON(http.StatusOK, cou.Title)
}

func deleteCourse(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	stmt, err := db.Prepare("DELETE FROM course where id=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	stmt.Exec(id)

	return c.NoContent(http.StatusOK)
}
