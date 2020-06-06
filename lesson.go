package main

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

// Lesson contains info about the lesson of particular course
type Lesson struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	ModuleID int    `json:"moduleID"`
}

func createLesson(c echo.Context) error {
	l := &Lesson{}
	if err := c.Bind(l); err != nil {
		return err
	}

	insertLessonSQL := "INSERT INTO lesson(title, module_id) VALUES(?,?)"
	stmt, err := db.Prepare(insertLessonSQL)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(l.Title, l.ModuleID)

	return c.JSON(http.StatusCreated, l.Title)
}

func getLessons(c echo.Context) error {
	var ls []Lesson = make([]Lesson, 0)

	rows, err := db.Query("SELECT id, title, module_id FROM lesson")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		l := Lesson{}
		rows.Scan(&l.ID, &l.Title, &l.ModuleID)
		ls = append(ls, l)
	}

	return c.JSON(http.StatusOK, ls)
}

func getLessonByID(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	row := db.QueryRow("SELECT id, title, module_id FROM lesson WHERE id=?", id)
	l := Lesson{}
	row.Scan(&l.ID, &l.Title, &l.ModuleID)

	return c.JSON(http.StatusOK, l)
}

func updateLesson(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	l := &Lesson{}

	if err := c.Bind(l); err != nil {
		return err
	}

	stmt, err := db.Prepare("UPDATE lesson SET title=? WHERE id=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(l.Title, id)

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
