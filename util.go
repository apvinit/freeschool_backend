package main

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo"
)

func uploadMedia(c echo.Context) error {

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

	splitName := strings.Split(file.Filename, ".")

	fileName := uuid.New().String() + "." + splitName[1]

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

	return c.JSON(http.StatusOK, map[string]string{"id": fileName})
}

func deleteMedia(fileName string) error {
	err := os.Remove(filepath.Join("freeschool", "media", fileName))
	if err != nil {
		return err
	}
	return nil
}

func getMedia(c echo.Context) error {
	fileName := c.Param("file")
	return c.File("freeschool/media/" + fileName)
}
