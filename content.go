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
	ct := &Content{}
	if err := c.Bind(ct); err != nil {
		return err
	}

	insertContentSQL := "INSERT INTO content(title, description, data, lesson_ID) VALUES(?,?,?,?)"
	stmt, err := db.Prepare(insertContentSQL)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(ct.Title, ct.Description, ct.Data, ct.LessonID)

	return c.JSON(http.StatusCreated, ct)
}

func getContents(c echo.Context) error {
	var con []Content = make([]Content, 0)

	rows, err := db.Query("SELECT id, title, description, data, lesson_id FROM content")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		co := Content{}
		rows.Scan(&co.ID, &co.Title, &co.Description, &co.Data, &co.LessonID)
		con = append(con, co)
	}

	return c.JSON(http.StatusOK, con)
}

func getContentByID(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	row := db.QueryRow("SELECT id, title, description, data, lesson_id FROM content WHERE id = ?", id)
	co := Content{}
	row.Scan(&co.ID, &co.Title, &co.Description, &co.Data, &co.LessonID)

	return c.JSON(http.StatusOK, co)
}

func updateContent(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	co := &Content{}

	if err := c.Bind(co); err != nil {
		return err
	}

	stmt, err := db.Prepare("UPDATE content SET title=? WHERE id=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(co.Title, id)

	return c.JSON(http.StatusOK, co)
}

func deleteContent(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	stmt, err := db.Prepare("DELETE FROM content where id=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	stmt.Exec(id)

	return c.NoContent(http.StatusOK)

}
