package main

import (
	"net/http"
	"strings"

	"github.com/labstack/echo"
)

// Language for keeping lessons
type Language struct {
	Lang  string `json:"lang,omitempty"`
	Title string `json:"title,omitempty"`
}

func createLanguage(c echo.Context) error {
	l := &Language{}
	if err := c.Bind(l); err != nil {
		return err
	}

	l.Lang = strings.Trim(l.Lang, " ")

	insertLangSQL := "INSERT INTO languages(lang,title) VALUES(?,?)"
	stmt, err := db.Prepare(insertLangSQL)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(l.Lang, l.Title)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, l)
}

func getLanguages(c echo.Context) error {
	var ls []Language = make([]Language, 0)
	rows, err := db.Query("SELECT * FROM languages")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		l := Language{}
		rows.Scan(&l.Lang, &l.Title)
		ls = append(ls, l)
	}

	return c.JSON(http.StatusOK, ls)
}

func updateLanguage(c echo.Context) error {
	lang := c.Param("id")
	l := &Language{}

	if err := c.Bind(l); err != nil {
		return err
	}

	stmt, err := db.Prepare("UPDATE languages SET lang=?, title=? WHERE lang=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(l.Lang, l.Title, lang)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusAccepted)
}

func deleteLanguage(c echo.Context) error {
	lang := c.Param("id")
	_, err := db.Exec("DELETE FROM languages WHERE lang = ?", lang)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
