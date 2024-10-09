package main

import (
	// built-ins
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	// 3rd party
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// type user struct {
// 	UserID    string     `json:"booklistID"`
// 	Email     string     `json:"email"`
// 	Username  string     `json:"username"`
// 	Password  string     `json:"password"`
// }

type bookList struct {
	BooklistID   uint      `json:"booklistID" db:"booklist_id,omitempty"`
	BooklistName string    `json:"booklistName" db:"booklist_name"`
	DTMCreated   time.Time `json:"dtmcreated" db:"booklist_name"`
	// BooklistBooks []book    `json:"booklistBooks"`
}

type book struct {
	BookID      uint      `json:"bookID"`
	BookTitle   string    `json:"bookTitle"`
	DTMAdded    time.Time `json:"dtmadded"`
	BookAuthors []author  `json:"bookAuthors"`
}

type author struct {
	AuthorID   uint   `json:"authorID"`
	BookID     string `json:"bookID"`
	AuthorName string `json:"authorName"`
}

// GET Handlers

// GET all bookLists for user
func getBookLists(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		user_id := 1
		var bookLists []bookList

		rows, err := db.Query("SELECT booklist_id, booklist_name, dtm_created FROM booklists WHERE user_id = $1", user_id)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		for rows.Next() {
			var bookList bookList
			if err := rows.Scan(&bookList.BooklistID, &bookList.BooklistName, &bookList.DTMCreated); err != nil {
				log.Fatal(err)
			}
			bookLists = append(bookLists, bookList)
		}

		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, bookLists)
	}
}

// get any booklist by ID
// func bookListId(id string) (*bookList, error) {
// 	for i, b := range bookLists {
// 		if b.BooklistID == id {
// 			return &bookLists[i], nil
// 		}
// 	}
// 	return nil, errors.New("booklist not found")
// }

// GET a single booklist and its books by ID
// func getBookListById(c *gin.Context) {
// 	id := c.Param("id")

// 	// TODO: replace with using a sql query
// 	// bookList, err := bookListId(id)

// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"message": "booklist not found"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, bookList)
// }

// POST Handlers
func createBookList(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		var newBookList bookList
		user_id := 1

		// bind JSON request body to bookList
		if err := c.BindJSON(&newBookList); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var booklist_id uint
		// TODO: limit individual users from creating booklists with the same name
		err := db.QueryRow("INSERT INTO booklists (booklist_name, dtm_created, user_id) SELECT $1, $2, $3 WHERE NOT EXISTS (SELECT user_id FROM booklists WHERE user_id = 1 AND booklist_name = $1) RETURNING booklist_id", newBookList.BooklistName, newBookList.DTMCreated, user_id).Scan(&booklist_id)

		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusConflict, gin.H{"message": "booklist name already exists for user_id"})
				return
			} else {
				log.Fatal(err)
			}
		}

		newBookList.BooklistID = booklist_id

		c.JSON(http.StatusCreated, newBookList)
	}
}

// add a book to a booklist based on the user ID and booklist ID
// func createNewBook(c *gin.Context) {
// 	// create the book variable that will be added to the booklist
// 	var newBook book

// 	// get the booklist to add the book to by id first
// 	id := c.Param("id")
// 	bookList, err := bookListId(id)
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"message": "booklist not found"})
// 		return
// 	}

// 	if err := c.BindJSON(&newBook); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// bookList.BooklistBooks = append(bookList.BooklistBooks, newBook)
// 	c.JSON(http.StatusCreated, bookList)
// }

// connect to the database and return a DB object
func dbConnect() *sql.DB {
	connStr := os.Getenv("BOOKISH_DB_CONN_STR")
	// connStr := "host=localhost port=5432 user=shantel password=password dbname=bookish sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func main() {
	// creating the DB object to be used by other functions
	db := dbConnect()
	// close the db connections if the program ends
	defer db.Close()

	router := gin.Default()
	router.GET("/booklists", getBookLists(db))
	// router.GET("/booklists/:id", getBookListById)
	router.POST("/booklists", createBookList(db))
	// router.POST("/booklists/create/:id", createNewBook)
	router.Run("localhost:8080")
}
