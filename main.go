package main

import (
	"net/http"

	"errors"

	"github.com/gin-gonic/gin"
)

// The json fields convert the struct name from capital letters to lower case
// so it can be worked as JSON data in the API
type book struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}

var books = []book{
	{ID: "1", Title: "In Search of Lost Time", Author: "Marcel Proust", Quantity: 2},
	{ID: "2", Title: "The Great Gatsby", Author: "F. Scott Fitzgerald", Quantity: 5},
	{ID: "3", Title: "War and Peace", Author: "Leo Tolstoy", Quantity: 6},
}

func getBooks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, books)
}

func bookById(c *gin.Context) {
	id := c.Param("id")
	book, err := getBookById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}

	c.IndentedJSON(http.StatusOK, book)
}

func checkoutBook(c *gin.Context) {
	// query parameter of "id", check if the id is ok or not
	id, ok := c.GetQuery("id")

	// if ok == false, return error
	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter."})
		return
	}

	// if ok, return book
	book, err := getBookById(id)

	// if error is not nil, return error msg
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}

	// if the book's quantity = 0, return a message
	if book.Quantity <= 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Book not available."})
		return
	}

	// if book's quantity can be minus 1 (ie checkout), return the book
	book.Quantity -= 1
	c.IndentedJSON(http.StatusOK, book)
}

func returnBook(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter."})
		return
	}

	book, err := getBookById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}

	book.Quantity += 1
	c.IndentedJSON(http.StatusOK, book)
}

func getBookById(id string) (*book, error) {
	for i, b := range books {
		if b.ID == id {
			return &books[i], nil
			// pointer is used to modify the attribute of the book,
			// field of the struct from different function
		}
	}

	return nil, errors.New("book not found")
}

/*
	1. We created a var newBook that is a type of book
	2. bind the JSON of the requested data to newBook.
	3. if there is an error, it will bypass the function via return
	4. if no error, we will append/bind the book to the newBook we are creating
	5. Return the book we just created with a status code of statusCreated
*/
func createBook(c *gin.Context) {
	var newBook book

	if err := c.BindJSON(&newBook); err != nil {
		return
	}

	books = append(books, newBook)
	c.IndentedJSON(http.StatusCreated, newBook)
}

/*
1. router variable
2. router comes from Gin to handle different route
3. handle /books route
4. and if you go to localhost:8080, you will get the books
5. the books will take in a Gin Context
6. all the information about the request and return a response
7. IndentedJSON format the data into workable JSON data
8. if statusOK, send in the data "books"
9. since there is a struct "book", the data will be serialised as JSON object
*/
func main() {
	router := gin.Default()
	router.GET("/books", getBooks)
	router.GET("/books/:id", bookById)
	router.POST("/books", createBook)
	// updating something exist = PATCH
	router.PATCH("/checkout", checkoutBook)
	router.PATCH("/return", returnBook)
	router.Run("localhost:8080")
}
