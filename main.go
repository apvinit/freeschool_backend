package main

import (
	"database/sql"
	"log"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/mattn/go-sqlite3"
)

var topics map[int]*Topic = make(map[int]*Topic, 0)
var lessons map[int]*Lesson = make(map[int]*Lesson, 0)
var contents map[int]*Content = make(map[int]*Content, 0)

var db *sql.DB

func initDB(db *sql.DB) {
	createCategoryTable := `
		CREATE TABLE IF NOT EXISTS category (
			"id" INTEGER PRIMARY KEY AUTOINCREMENT,
			"title" TEXT,
			"description" TEXT
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
}

func main() {

	db, _ = sql.Open("sqlite3", "./backend.db")
	// if err != nil {
	// 	panic(err)
	// }
	defer db.Close()

	initDB(db)

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

	api.POST("/topics", createTopic)
	api.GET("/topics", getTopics)
	api.GET("/topics/:id", getTopicByID)
	api.PUT("/topics/:id", updateTopic)
	api.DELETE("/topics/:id", deleteTopic)

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

	e.Start(":8888")

}
