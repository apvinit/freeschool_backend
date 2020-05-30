package main

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

// Topic is a slice of lessons with some meta data
type Topic struct {
	ID          int      `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Lessons     []Lesson `json:"lessons"`
}

func createTopic(c echo.Context) error {
	tp := &Topic{}
	if err := c.Bind(tp); err != nil {
		return err
	}
	topics[tp.ID] = tp
	return c.JSON(http.StatusCreated, tp)
}

func getTopics(c echo.Context) error {
	var tp []Topic = make([]Topic, 0)

	for _, top := range topics {
		tp = append(tp, *top)
	}
	return c.JSON(http.StatusOK, tp)
}

func getTopicByID(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	return c.JSON(http.StatusOK, topics[id])
}

func updateTopic(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	top := &Topic{}

	if err := c.Bind(&top); err != nil {
		return err
	}

	topics[id] = top

	return c.JSON(http.StatusOK, top)
}

func deleteTopic(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	delete(topics, id)

	return c.NoContent(http.StatusOK)

}
