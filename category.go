package main

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

// Category contains the info about the category of the courses
type Category struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	IconURL string `json:"iconURL,omitempty"`
}

func createCategory(c echo.Context) error {

	cat := &Category{}
	if err := c.Bind(cat); err != nil {
		return err
	}
	categories[cat.ID] = cat
	return c.JSON(http.StatusCreated, cat)
}

func getCategories(c echo.Context) error {
	var cat []Category = make([]Category, 0)

	for _, ct := range categories {
		cat = append(cat, *ct)
	}
	return c.JSON(http.StatusOK, cat)
}

func getCategoryByID(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	return c.JSON(http.StatusOK, categories[id])
}

func updateCategory(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	cat := &Category{}

	if err := c.Bind(&cat); err != nil {
		return err
	}

	categories[id] = cat

	return c.JSON(http.StatusOK, cat)
}

func deleteCategory(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	delete(categories, id)

	return c.NoContent(http.StatusOK)
}
