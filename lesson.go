package main

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

// Lesson contains info about the lesson of particular course
type Lesson struct {
	ID          int    `json:"id,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	ModuleID    int    `json:"module_id,omitempty"`
	Draft       bool   `json:"draft,omitempty"`
}

func createLesson(c echo.Context) error {
	l := &Lesson{}
	if err := c.Bind(l); err != nil {
		return err
	}

	insertLessonSQL := "INSERT INTO lesson(title, description, module_id, draft) VALUES(?,?,?,?)"
	stmt, err := db.Prepare(insertLessonSQL)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(l.Title, l.Description, l.ModuleID, l.Draft)

	return c.JSON(http.StatusCreated, l)
}

func getLessons(c echo.Context) error {
	var ls []Lesson = make([]Lesson, 0)

	lid := c.QueryParam("module_id")
	if len(lid) != 0 {
		moduleID, err := strconv.Atoi(lid)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "malformatted module_id"})
		}

		rows, err := db.Query("SELECT id, title, description, module_id, draft FROM lesson WHERE module_id=?", moduleID)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			l := Lesson{}
			rows.Scan(&l.ID, &l.Title, &l.Description, &l.ModuleID, &l.Draft)
			ls = append(ls, l)
		}

		return c.JSON(http.StatusOK, ls)
	}

	rows, err := db.Query("SELECT id, title, description, module_id, draft FROM lesson")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		l := Lesson{}
		rows.Scan(&l.ID, &l.Title, &l.Description, &l.ModuleID, &l.Draft)
		ls = append(ls, l)
	}

	return c.JSON(http.StatusOK, ls)
}

func getLessonByID(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	row := db.QueryRow("SELECT id, title, description, module_id, draft FROM lesson WHERE id=?", id)
	l := Lesson{}
	row.Scan(&l.ID, &l.Title, &l.Description, &l.ModuleID, &l.Draft)

	return c.JSON(http.StatusOK, l)
}

func updateLesson(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	l := &Lesson{}

	if err := c.Bind(l); err != nil {
		return err
	}

	stmt, err := db.Prepare("UPDATE lesson SET title=?, description=?, module_id=?, draft=? WHERE id=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(l.Title, l.Description, l.ModuleID, l.Draft, id)

	return c.JSON(http.StatusOK, l)
}

func deleteLesson(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	stmt, err := db.Prepare("DELETE FROM lesson where id=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	stmt.Exec(id)

	return c.NoContent(http.StatusOK)

}
