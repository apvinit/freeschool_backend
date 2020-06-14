package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo"
)

// Content is the core element. It has the actual course content.
type Content struct {
	ID          int    `json:"id,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	LessonID    int    `json:"lesson_id,omitempty"`
	ContentType string `json:"content_type,omitemtpy"`
	Data        string `json:"data,omitempty"`
	Draft       bool   `json:"draft,omitempty"`
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

	cid := c.QueryParam("lesson_id")
	if len(cid) != 0 {
		lessonID, err := strconv.Atoi(cid)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "malformatted module_id"})
		}

		rows, err := db.Query("SELECT id, title, description, data FROM content WHERE lesson_id= ?", lessonID)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			co := Content{}
			rows.Scan(&co.ID, &co.Title, &co.Description, &co.Data)
			con = append(con, co)
		}

		return c.JSON(http.StatusOK, con)
	}

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

	row := db.QueryRow("SELECT data FROM content WHERE id = ?", id)
	var dataID string
	row.Scan(&dataID)

	deleteMedia(dataID)
	deleteTransacodedMedia(dataID)

	stmt, err := db.Prepare("DELETE FROM content where id=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	stmt.Exec(id)

	return c.NoContent(http.StatusOK)

}

func uploadContent(c echo.Context) error {

	// Source
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	fileName := uuid.New().String()

	// Destination
	dst, err := os.Create(filepath.Join("media", fileName))
	if err != nil {
		return err
	}

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}
	dst.Close()

	// Transcode the uploaded file to HLS
	go transcodeToHLS(fileName)

	return c.JSON(http.StatusOK,
		map[string]string{"id": fileName})
}

func streamFileSegment(c echo.Context) error {
	name := c.Param("file")
	segment := c.Param("segment")
	return c.File("transcoded/" + name + "/" + segment)
}

func deleteTransacodedMedia(fileName string) {
	err := os.RemoveAll(filepath.Join("transcoded", fileName))
	if err != nil {
		fmt.Println("Error Deleting directory", err.Error())
		return
	}
	fmt.Println("File Deleted Successfully")
}
