package main

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

// Module is a slice of lessons with some meta data
type Module struct {
	ID          int    `json:"id,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	CourseID    int    `json:"courseID,omitempty"`
}

func createModule(c echo.Context) error {
	m := &Module{}
	if err := c.Bind(m); err != nil {
		return err
	}

	insertCourseSQL := "INSERT INTO module(title, course_id) VALUES(?,?)"

	stmt, err := db.Prepare(insertCourseSQL)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(m.Title, m.CourseID)

	return c.JSON(http.StatusCreated, m)
}

func getModules(c echo.Context) error {
	var mod []Module = make([]Module, 0)

	mid := c.QueryParam("course_id")
	if len(mid) != 0 {
		courseID, err := strconv.Atoi(mid)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "malformatted course_id"})
		}

		rows, err := db.Query("SELECT id, title FROM module WHERE course_id=?", courseID)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			m := Module{}
			rows.Scan(&m.ID, &m.Title)
			mod = append(mod, m)
		}

		return c.JSON(http.StatusOK, mod)

	}

	row, err := db.Query("SELECT id, title, course_id FROM module")
	if err != nil {
		return err
	}
	defer row.Close()

	for row.Next() {
		m := Module{}
		row.Scan(&m.ID, &m.Title, &m.CourseID)
		mod = append(mod, m)
	}
	return c.JSON(http.StatusOK, mod)
}

func getModuleByID(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	row := db.QueryRow("SELECT id, title FROM module WHERE id=?", id)
	m := Module{}
	row.Scan(&m.ID, &m.Title)

	return c.JSON(http.StatusOK, m)
}

func updateModule(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	m := &Module{}

	if err := c.Bind(&m); err != nil {
		return err
	}

	stmt, err := db.Prepare("UPDATE module SET title=? WHERE id=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(m.Title, id)

	return c.JSON(http.StatusOK, m.Title)
}

func deleteModule(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	stmt, err := db.Prepare("DELETE FROM module where id=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	stmt.Exec(id)

	return c.NoContent(http.StatusOK)

}
