package main

import (
	"net/http"

	"github.com/labstack/echo"
)

// Tag for keeping lessons
type Tag struct {
	ID    int    `json:"id,omitempty"`
	Title string `json:"title,omitempty"`
}

func createTag(c echo.Context) error {
	t := &Tag{}
	if err := c.Bind(t); err != nil {
		return err
	}

	insertTagSQL := "INSERT INTO tags(title) VALUES(?)"
	stmt, err := db.Prepare(insertTagSQL)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(t.Title)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, t)
}

func getTags(c echo.Context) error {
	var ts []Tag = make([]Tag, 0)
	rows, err := db.Query("SELECT * FROM tags")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		t := Tag{}
		rows.Scan(&t.ID, &t.Title)
		ts = append(ts, t)
	}

	return c.JSON(http.StatusOK, ts)
}

func updateTag(c echo.Context) error {
	id := c.Param("id")
	t := &Tag{}

	if err := c.Bind(t); err != nil {
		return err
	}

	stmt, err := db.Prepare("UPDATE tags SET title=? WHERE id=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(t.Title, id)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusAccepted)
}

func deleteTag(c echo.Context) error {
	id := c.Param("id")
	_, err := db.Exec("DELETE FROM tags WHERE id = ?", id)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
