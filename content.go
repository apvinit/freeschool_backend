package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

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

	insertContentSQL := "INSERT INTO contents(title, description, lesson_ID, content_type, data, draft) VALUES(?,?,?,?,?,?)"
	stmt, err := db.Prepare(insertContentSQL)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(ct.Title, ct.Description, ct.LessonID, ct.ContentType, ct.Data, ct.Draft)

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

		rows, err := db.Query("SELECT id, title, description, lesson_id, content_type, data, draft FROM contents WHERE lesson_id= ?", lessonID)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			co := Content{}
			rows.Scan(&co.ID, &co.Title, &co.Description, &co.LessonID, &co.ContentType, &co.Data, &co.Draft)
			con = append(con, co)
		}

		return c.JSON(http.StatusOK, con)
	}

	rows, err := db.Query("SELECT id, title, description, lesson_id, content_type, data, draft FROM contents")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		co := Content{}
		rows.Scan(&co.ID, &co.Title, &co.Description, &co.LessonID, &co.ContentType, &co.Data, &co.Draft)
		con = append(con, co)
	}

	return c.JSON(http.StatusOK, con)
}

func getContentsForLesson(lessonID int) []Content {
	var con []Content = make([]Content, 0)
	rows, err := db.Query("SELECT id, title, description, lesson_id, content_type, data, draft FROM contents WHERE lesson_id= ?", lessonID)
	if err != nil {
		return []Content{}
	}
	defer rows.Close()

	for rows.Next() {
		co := Content{}
		rows.Scan(&co.ID, &co.Title, &co.Description, &co.LessonID, &co.ContentType, &co.Data, &co.Draft)
		con = append(con, co)
	}
	return con
}

func getContentByID(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	row := db.QueryRow("SELECT id, title, description, lesson_id, content_type, data, draft FROM contents WHERE id = ?", id)
	co := Content{}
	row.Scan(&co.ID, &co.Title, &co.Description, &co.LessonID, &co.ContentType, &co.Data, &co.Draft)

	return c.JSON(http.StatusOK, co)
}

func updateContent(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	co := &Content{}

	if err := c.Bind(co); err != nil {
		return err
	}

	stmt, err := db.Prepare("UPDATE contents SET title=?, description=?, lesson_id=?, content_type=?, data=?, draft=? WHERE id=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(co.Title, co.Description, co.LessonID, co.ContentType, co.Data, co.Draft, id)

	return c.JSON(http.StatusOK, co)
}

func deleteContent(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	row := db.QueryRow("SELECT data FROM contents WHERE id = ?", id)
	var dataID string
	row.Scan(&dataID)

	deleteMedia(dataID)
	deleteTransacodedMedia(dataID)

	stmt, err := db.Prepare("DELETE FROM contents where id=?")
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

	ext := strings.Split(file.Filename, ".")[1]
	fileName := uuid.New().String() + "." + ext

	// Destination
	dst, err := os.Create(filepath.Join("freeschool", "media", fileName))
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
	return c.File("freeschool/transcoded/" + name + "/" + segment)
}

func deleteTransacodedMedia(fileName string) {
	err := os.RemoveAll(filepath.Join("freeschool", "transcoded", fileName))
	if err != nil {
		fmt.Println("Error Deleting directory", err.Error())
		return
	}
	fmt.Println("File Deleted Successfully")
}
