package main

// This is all the work of Gebek studying Golang at KBTU

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Struct for data structure
type Item struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

// Temp DATA storage
var items = []Item{
	{ID: 1, Name: "Coca Cola", Price: 450},
	{ID: 2, Name: "Uzbek samsa", Price: 400},
	{ID: 3, Name: "Bahandi Burger", Price: 1150},
}

// CRUD

func getItems(c *gin.Context) {
	c.JSON(http.StatusOK, items)
}

func createItem(c *gin.Context) {
	var newItem Item

	if err := c.BindJSON(&newItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error occured": err.Error()})
		return
	}
	newItem.ID = len(items) + 1
	items = append(items, newItem)
	c.JSON(http.StatusCreated, newItem)
}

func updateItem(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var updatedItem Item

	if err := c.BindJSON(&updatedItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error occured": err.Error()})
		return
	}

	for i, item := range items {
		if item.ID == id {
			items[i] = updatedItem
			c.JSON(http.StatusOK, updatedItem)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error occured": "this item is not found"})
}

func deleteItem(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	for i, item := range items {
		if item.ID == id {
			items = append(items[:i], items[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"Success": "item is deleted"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error occured": "this item is not found"})
}

func main() {
	fmt.Println("This is the work of Gebek")

	r := gin.Default()

	config := cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}

	r.Use(cors.New(config))

	// Register and login for auth
	r.POST("/register", Register)
	r.POST("/login", Login)

	// Routes

	authorized := r.Group("/")
	authorized.Use(AuthMiddleware())
	{
		// User
		authorized.GET("/items", getItems)

		// Admin protected
		authorized.POST("/items", RoleCheckMiddleware("admin"), createItem)
		authorized.PUT("/items/:id", RoleCheckMiddleware("admin"), updateItem)
		authorized.DELETE("/items/:id", RoleCheckMiddleware("admin"), deleteItem)
	}

	r.Run(":8080")
}
