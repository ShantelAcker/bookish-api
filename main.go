package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type bookList struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	DTMCreated time.Time `json:"dtmcreated"`
}

var bookLists = []bookList{
	{ID: "1", Name: "Want to Read", DTMCreated: time.Now()},
	{ID: "2", Name: "Reading", DTMCreated: time.Now()},
	{ID: "3", Name: "Finished", DTMCreated: time.Now()},
}

// GET Handlers
func getBookLists(c *gin.Context) {
	c.JSON(http.StatusOK, bookLists)
}

// middleware for getting any book's ID
func bookId(id string) (*bookList, error) {
	for i, b := range bookLists {
		if b.ID == id {
			return &bookLists[i], nil
		}
	}
	return nil, errors.New("book not found")
}

func getBookListById(c *gin.Context) {
	id := c.Param("id")
	bookList, err := bookId(id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "booklist not found"})
	}

	c.JSON(http.StatusOK, bookList)
}

// POST Handlers
func createBookList(c *gin.Context) {

}

func main() {
	router := gin.Default()
	router.GET("/booklists", getBookLists)
	router.GET("/booklists/:id", getBookListById)
	router.Run("localhost:8080")
}
