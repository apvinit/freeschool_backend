package main

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

// Category contains the info about the category of the courses
type Category struct {
	ID    int    `json:"id,omitempty"`
	Title string `json:"title,omitempty"`
	Cover string `json:"cover,omitempty"`
	Lang  string `json:"lang,omitempty"`
}

func createCategory(c echo.Context) error {

	cat := &Category{}
	if err := c.Bind(cat); err != nil {
		return err
	}

	insertCategorySQL := "INSERT INTO category(title, cover, lang) VALUES (?,?,?)"

	stmt, err := db.Prepare(insertCategorySQL)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(cat.Title, cat.Cover, cat.Lang)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, map[string]string{"status": "created"})
}

func getCategories(c echo.Context) error {

	row, err := db.Query("SELECT id, title FROM category")
	if err != nil {
		return err
	}
	defer row.Close()

	var cat []Category = make([]Category, 0)

	for row.Next() {
		ct := Category{}
		err = row.Scan(&ct.ID, &ct.Title)
		if err != nil {
			return err
		}
		cat = append(cat, ct)
	}
	return c.JSON(http.StatusOK, cat)
}

func getCategoryByID(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	row := db.QueryRow("SELECT id, title FROM category WHERE id=?", id)

	cat := Category{}

	row.Scan(&cat.ID, &cat.Title)

	return c.JSON(http.StatusOK, cat)
}

func updateCategory(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	cat := &Category{}

	if err := c.Bind(&cat); err != nil {
		return err
	}

	stmt, err := db.Prepare("UPDATE category SET title=? WHERE id=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(cat.Title, id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, cat.Title)
}

func deleteCategory(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	stmt, err := db.Prepare("delete from category where id=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	stmt.Exec(id)

	return c.NoContent(http.StatusOK)
}
