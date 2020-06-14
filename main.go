package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func initDB(db *sql.DB) {
	createCategoryTable := `
		CREATE TABLE IF NOT EXISTS category (
			"id" INTEGER PRIMARY KEY AUTOINCREMENT,
			"title" TEXT,
			"cover" TEXT,
			"lang" TEXT
		);
	`

	statement, err := db.Prepare(createCategoryTable)
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()
	statement.Close()

	createCourseTable := `
		CREATE TABLE IF NOT EXISTS course(
			"id" INTEGER PRIMARY KEY AUTOINCREMENT,
			"title" TEXT,
			"description" TEXT,
			"category_id" INTEGER
		);	
	`
	statement, err = db.Prepare(createCourseTable)
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()
	statement.Close()

	createModuleTable := `
		CREATE TABLE IF NOT EXISTS module(
			"id" INTEGER PRIMARY KEY AUTOINCREMENT,
			"title" TEXT,
			"description" TEXT,
			"course_id" INTEGER
		);
	`
	statement, err = db.Prepare(createModuleTable)
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()
	statement.Close()

	createLessonTable := `
		CREATE TABLE IF NOT EXISTS lesson(
			"id" INTEGER PRIMARY KEY AUTOINCREMENT,
			"title" TEXT,
			"module_id" TEXT
		);
	`
	statement, err = db.Prepare(createLessonTable)
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()
	statement.Close()

	createContentTable := `
		CREATE TABLE IF NOT EXISTS content(
			"id" INTEGER PRIMARY KEY AUTOINCREMENT,
			"title" TEXT,
			"description" TEXT,
			"lesson_id" INTEGER,
			"content_type" TEXT,
			"data" TEXT
		);
	`

	statement, err = db.Prepare(createContentTable)
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()
	statement.Close()

}

func main() {

	db, _ = sql.Open("sqlite3", "./backend.db")
	// if err != nil {
	// 	panic(err)
	// }
	defer db.Close()

	initDB(db)

	// setup upload directories
	if _, err := os.Stat("media"); os.IsNotExist(err) {
		os.Mkdir("media", 0755)
	}

	if _, err := os.Stat("transcoded"); os.IsNotExist(err) {
		os.Mkdir("transcoded", 0755)
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.CORS())
	api := e.Group("api")

	api.POST("/categories", createCategory)
	api.GET("/categories", getCategories)
	api.GET("/categories/:id", getCategoryByID)
	api.PUT("/categories/:id", updateCategory)
	api.DELETE("/categories/:id", deleteCategory)

	api.POST("/courses", createCourse)
	api.GET("/courses", getCourses)
	api.GET("/courses/:id", getCourseByID)
	api.PUT("/courses/:id", updateCourse)
	api.DELETE("/courses/:id", deleteCourse)

	api.POST("/modules", createModule)
	api.GET("/modules", getModules)
	api.GET("/modules/:id", getModuleByID)
	api.PUT("/modules/:id", updateModule)
	api.DELETE("/modules/:id", deleteModule)

	api.POST("/lessons", createLesson)
	api.GET("/lessons", getLessons)
	api.GET("/lessons/:id", getLessonByID)
	api.PUT("/lessons/:id", updateLesson)
	api.DELETE("/lessons/:id", deleteLesson)

	api.POST("/contents", createContent)
	api.GET("/contents", getContents)
	api.GET("/contents/:id", getContentByID)
	api.PUT("/contents/:id", updateContent)
	api.DELETE("/contents/:id", deleteContent)
	api.POST("/contents/upload", uploadContent)
	api.GET("/contents/stream/:file/:segment", streamFileSegment)

	api.POST("/upload", uploadMedia)

	e.Start(":8888")

}
