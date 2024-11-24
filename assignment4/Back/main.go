package main

// This is all the work of Gebek studying Golang at KBTU

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

var validate = validator.New()

// prometheus
var requestsTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests",
	},
	[]string{"method", "status"},
)

var requestDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Histogram of HTTP request durations.",
		Buckets: prometheus.DefBuckets,
	},
	[]string{"method", "status"},
)

func init() {
	prometheus.MustRegister(requestDuration)
	prometheus.MustRegister(requestsTotal)
}

// Application metrics
var errorCount = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_errors_total",
		Help: "Total number of HTTP errors (4xx and 5xx)",
	},
	[]string{"method", "status"},
)

// Logging
var logger = logrus.New()

func logRequest(c *gin.Context) {
	start := time.Now()
	defer func() {
		duration := time.Since(start).Seconds()
		status := fmt.Sprintf("%d", c.Writer.Status())
		requestDuration.WithLabelValues(c.Request.Method, status).Observe(duration)

		// Log errors (4xx or 5xx status codes)
		if c.Writer.Status() >= 400 {
			errorCount.WithLabelValues(c.Request.Method, status).Inc()
		}
	}()

	c.Next()
}

// Struct for data structure
type Item struct {
	ID    int    `json:"id"`
	Name  string `json:"name" validate:"required,min=3"`
	Price int    `json:"price" validate:"required,gt=0"`
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
	logger.Info("Fetching items")
}

func createItem(c *gin.Context) {
	var newItem Item

	if err := c.BindJSON(&newItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	if err := validate.Struct(newItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		logger.Error(err.Error())
		return
	}

	newItem.ID = len(items) + 1
	items = append(items, newItem)
	c.JSON(http.StatusCreated, newItem)
	logger.Info("Item created successfully", logrus.Fields{"item_id": newItem.ID, "name": newItem.Name})
}

func updateItem(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		logger.Error("Failed to update item", logrus.Fields{"error": err.Error()})
		return
	}

	var updatedItem Item

	if err := c.BindJSON(&updatedItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error occured": err.Error()})
		logger.Error("Failed to update item", logrus.Fields{"error": err.Error()})
		return
	}

	for i, item := range items {
		if item.ID == id {
			items[i] = updatedItem
			c.JSON(http.StatusOK, updatedItem)
			logger.Info("Updating items")
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
		logger.Warn("Attempted to delete a non-existing item", logrus.Fields{"item_id": id})
		return
	}

	for i, item := range items {
		if item.ID == id {
			logger.Warn("Deleting item?")
			logger.Info(item)
			items = append(items[:i], items[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"Success": "item is deleted"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error occured": "this item is not found"})
	logger.Info("Deleting items")
}

// var csrfKey = []byte("a21170bc6a284405c068f5d8bef1198647023024c201f00da4303664b7d72fcf")

func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Content Security Policy: Restricts the sources from which content can be loaded.
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self'; object-src 'none';")

		// X-Frame-Options: Prevents the site from being embedded in a frame or iframe.
		c.Header("X-Frame-Options", "DENY")

		// X-Content-Type-Options: Prevents browsers from interpreting files as a different MIME type.
		c.Header("X-Content-Type-Options", "nosniff")

		// Strict-Transport-Security: Forces browsers to use HTTPS for all future requests.
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		// Referrer-Policy: Controls how much referrer information is sent with requests.
		c.Header("Referrer-Policy", "no-referrer")

		// Feature-Policy (or Permissions-Policy): Limits the use of certain features like geolocation, camera, etc.
		c.Header("Permissions-Policy", "geolocation=(self), microphone=()")

		// Set the X-XSS-Protection header to prevent some types of Cross-Site Scripting (XSS) attacks
		c.Header("X-XSS-Protection", "1; mode=block")

		// Allow-Control-Allow-Origin: Prevents cross-origin requests that are not allowed.
		c.Header("Access-Control-Allow-Origin", "*")

		// Proceed with the request
		c.Next()
	}
}

func main() {
	fmt.Println("This is the work of Gebek")

	r := gin.Default()

	logger.SetFormatter(&logrus.JSONFormatter{})

	config := cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}

	r.Use(SecurityHeadersMiddleware(), cors.New(config), logRequest)

	// csrfMiddleware := csrf.Protect(csrfKey, csrf.Secure(false))

	// r.Use(func(c *gin.Context) {
	// 	c.Header("X-CSRF-Token", csrf.Token(c.Request))
	// 	c.Next()
	// })
	// Register and login for auth
	r.POST("/register", Register)
	r.POST("/login", Login)

	// Routes
	authorized := r.Group("/")
	authorized.Use(AuthMiddleware())
	{
		// User routes
		authorized.GET("/items", getItems)

		// Admin-protected routes
		authorized.POST("/items", RoleCheckMiddleware("admin"), createItem)
		authorized.PUT("/items/:id", RoleCheckMiddleware("admin"), updateItem)
		authorized.DELETE("/items/:id", RoleCheckMiddleware("admin"), deleteItem)
		authorized.GET("/metrics", RoleCheckMiddleware("admin"), gin.WrapH(promhttp.Handler()))
	}

	// Run the server with TLS
	// err := r.RunTLS(":8080", "D:/Class/GO/assignment4/Back/cert.pem", "D:/Class/GO/assignment4/Back/key.pem")
	// if err != nil {
	// 	fmt.Println("Error starting server:", err)
	// }

	r.Run(":8080")
}
