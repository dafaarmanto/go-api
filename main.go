package main

import (
	"net/http"

	"errors"

	"github.com/gin-gonic/gin"
)

// Convert to JSON
type book struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}

// Create book slices
var books = []book{
	{ID: "1", Title: "Book 1", Author: "Author 1", Quantity: 2},
	{ID: "2", Title: "Book 2", Author: "Author 2", Quantity: 5},
	{ID: "3", Title: "Book 3", Author: "Author 3", Quantity: 6},
}

func getBooks(c *gin.Context) {
	// gin.Context is a struct that contains information about the request
	// IndentedJSON is a function that converts the data to JSON in a pretty format
	// http.StatusOK is a constant that returns the status code 200 (OK)
	// books is the slice of books to be converted to JSON
	c.IndentedJSON(http.StatusOK, books) // Return the books
}

func createBook(c *gin.Context) {
	var newBook book                             // Get the data from the request
	if err := c.BindJSON(&newBook); err != nil { // Bind the data to the book struct
		return
	}

	books = append(books, newBook)              // Append the new book to the books slice
	c.IndentedJSON(http.StatusCreated, newBook) // Return the new book
}

func bookById(c *gin.Context) {
	id := c.Param("id")          // Get the id from the URL
	book, err := getBookById(id) // Get the book by id
	if err != nil {              // Return an error if the book is not found
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()}) // Return an error if the book is not found
		return
	}
	c.IndentedJSON(http.StatusOK, book) // Return the book
}

func checkoutBook(c *gin.Context) {
	id, ok := c.GetQuery("id") // Get the id from the query

	if !ok { // Check if the id is not empty
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Missing id"}) // Return an error if the id is not provided
		return
	}

	book, err := getBookById(id) // Get the book by id

	if err != nil { // Return an error if the book is not found
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()}) // Return an error if the book is not found
		return
	}

	if book.Quantity == 0 { // Check if the book is out of stock
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Book out of stock"}) // Return an error if the book is out of stock
		return
	}

	book.Quantity -= 1                  // Decrease the quantity of the book by 1
	c.IndentedJSON(http.StatusOK, book) // Return the book
}

func returnBook(c *gin.Context) {
	id, ok := c.GetQuery("id") // Get the id from the query

	if !ok { // Check if the id is not empty
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Missing id"}) // Return an error if the id is not provided
		return
	}

	book, err := getBookById(id) // Get the book by id

	if err != nil { // Return an error if the book is not found
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()}) // Return an error if the book is not found
		return
	}

	book.Quantity += 1                  // Increase the quantity of the book by 1
	c.IndentedJSON(http.StatusOK, book) // Return the book
}

func getBookById(id string) (*book, error) {
	for i, b := range books { // Loop through the books
		if b.ID == id { // Check if the id matches
			return &books[i], nil // Return the book if it is found
		}
	}
	return nil, errors.New("Book not found") // Return an error if the book is not found
}

func main() {
	router := gin.Default()                 // Create a router from gin
	router.GET("/books", getBooks)          // Add a route to the router
	router.GET("/books/:id", bookById)      // Get a book by id
	router.POST("/books", createBook)       // Create a new book
	router.PATCH("/checkout", checkoutBook) // Checkout a book
	router.PATCH("/return", returnBook)     // Return a book
	router.Run("localhost:8080")            // Start the server on localhost:8080
}
