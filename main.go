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

// GET all bookLists
func getBookLists(c *gin.Context) {
	c.JSON(http.StatusOK, bookLists)
}

// middleware for getting any booklist ID
func bookListId(id string) (*bookList, error) {
	for i, b := range bookLists {
		if b.ID == id {
			return &bookLists[i], nil
		}
	}
	return nil, errors.New("booklist not found")
}

// GET a single booklist by ID
func getBookListById(c *gin.Context) {
	id := c.Param("id")
	bookList, err := bookListId(id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "booklist not found"})
	}

	c.JSON(http.StatusOK, bookList)
}

// POST Handlers
func createBookList(c *gin.Context) {
	var newBookList bookList

	// bind JSON response to bookList
	if err := c.BindJSON(&newBookList); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

	bookLists = append(bookLists, newBookList)
	c.JSON(http.StatusCreated, newBookList)
}

func main() {
	router := gin.Default()
	router.GET("/booklists", getBookLists)
	router.GET("/booklists/:id", getBookListById)
	router.POST("/booklists", createBookList)
	router.Run("localhost:8080")
}
