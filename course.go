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
	Cover       string `json:"cover,omitempty"`
	Lang        string `json:"lang,omitempty"`
	CategoryID  int    `json:"category_id,omitempty"`
	CreatedBy   string `json:"created_by,omitempty"`
	Draft       bool   `json:"draft,omitempty"`
}

func createCourse(c echo.Context) error {

	cou := &Course{}
	if err := c.Bind(cou); err != nil {
		return err
	}

	insertCourseSQL :=
		`INSERT INTO courses
			(title, description, category_id, cover, lang, created_by, draft) 
			VALUES(?,?,?,?,?,?,?)
	`

	stmt, err := db.Prepare(insertCourseSQL)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(cou.Title, cou.Description, cou.CategoryID, cou.Cover, cou.Lang, cou.CreatedBy, cou.Draft)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, cou)
}

func getCourses(c echo.Context) error {
	var cou []Course = make([]Course, 0)

	cid := c.QueryParam("category_id")
	if len(cid) != 0 {
		categoryID, err := strconv.Atoi(c.QueryParam("category_id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "malformatted category_id"})
		}

		rows, err := db.Query("SELECT id, title,description, category_id, cover, lang, created_by, draft FROM courses WHERE category_id=?", categoryID)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			co := Course{}
			rows.Scan(&co.ID, &co.Title, &co.Description, &co.CategoryID, &co.Cover, &co.Lang, &co.CreatedBy, &co.Draft)
			cou = append(cou, co)
		}

		return c.JSON(http.StatusOK, cou)
	}

	rows, err := db.Query("SELECT id, title,description, category_id, cover, lang, created_by, draft FROM courses")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		co := Course{}
		rows.Scan(&co.ID, &co.Title, &co.Description, &co.CategoryID, &co.Cover, &co.Lang, &co.CreatedBy, &co.Draft)
		cou = append(cou, co)
	}

	return c.JSON(http.StatusOK, cou)
}

func getCourseByID(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	row := db.QueryRow("SELECT id, title,description, category_id, cover, lang, created_by, draft FROM courses WHERE id=?", id)
	co := Course{}
	row.Scan(&co.ID, &co.Title, &co.Description, &co.CategoryID, &co.Cover, &co.Lang, &co.CreatedBy, &co.Draft)

	return c.JSON(http.StatusOK, co)
}

func updateCourse(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	cou := &Course{}

	if err := c.Bind(&cou); err != nil {
		return err
	}

	stmt, err := db.Prepare("UPDATE courses SET title=?,  description=?, category_id=?, cover=?, lang=?, created_by=?, draft=? WHERE id=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(cou.Title, cou.Description, cou.CategoryID, cou.Cover, cou.Lang, cou.CreatedBy, cou.Draft, id)

	return c.JSON(http.StatusOK, cou)
}

func deleteCourse(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	stmt, err := db.Prepare("DELETE FROM courses where id=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	stmt.Exec(id)

	return c.NoContent(http.StatusOK)
}
