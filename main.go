package main

import (
	// built-ins
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	// 3rd party
	"github.com/gin-gonic/gin"
)

type bookList struct {
	BooklistID    string    `json:"booklistID"`
	BooklistName  string    `json:"booklistName"`
	DTMCreated    time.Time `json:"dtmcreated"`
	BooklistBooks []book    `json:"booklistBooks"`
}

type book struct {
	BookID      string    `json:"bookID"`
	BookTitle   string    `json:"bookTitle"`
	DTMAdded    time.Time `json:"dtmadded"`
	BookAuthors []author  `json:"bookAuthors"`
}

type author struct {
	AuthorID   string `json:"authorID"`
	BookID     string `json:"bookID"`
	AuthorName string `json:"authorName"`
}

var bookLists = []bookList{}

// GET Handlers

// GET all bookLists
func getBookLists(c *gin.Context) {
	c.JSON(http.StatusOK, bookLists)
}

// middleware for getting any booklist ID
func bookListId(id string) (*bookList, error) {
	for i, b := range bookLists {
		if b.BooklistID == id {
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
	// reading in a dummy JSON file payload
	content, err := os.ReadFile("booklists.json")
	if err != nil {
		fmt.Println("unable to read json file")
	}
	err = json.Unmarshal(content, &bookLists)
	if err != nil {
		fmt.Println("unable to unmarshal json")
	}
	fmt.Println(bookLists[0].BooklistBooks[0].BookAuthors[0].AuthorName)

	router := gin.Default()
	router.GET("/booklists", getBookLists)
	router.GET("/booklists/:id", getBookListById)
	router.POST("/booklists", createBookList)
	router.Run("localhost:8080")
}
