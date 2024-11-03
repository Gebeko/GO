package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

type User struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"unique"`
	Age  int    `json:"age"`
}

type Task struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Users       []User `gorm:"many2many:task_users;"`
}

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "20120117"
	dbname   = "Midterm"
)

// CONNECT TO DB
func initDB() {
	var err error

	login := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err = gorm.Open(postgres.Open(login), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	fmt.Println("Success in connecting")

	db.AutoMigrate(&User{}, &Task{})
}

// Get all user with GORM
func GetUsersGORM(w http.ResponseWriter, r *http.Request) {
	var users []User
	query := db

	if age := r.URL.Query().Get("age"); age != "" {
		query = query.Where("age = ?", age)
	}

	if sort := r.URL.Query().Get("sort"); sort == "asc" {
		query = query.Order("name asc")
	} else if sort == "desc" {
		query = query.Order("name desc")
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			query = query.Limit(limit)
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			query = query.Offset(offset)
		}
	}

	if err := query.Find(&users).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(users)
}

// Create a new user
func CreateUserGORM(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := db.Create(&user).Error; err != nil {
		http.Error(w, "User with this name already exists", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Update an existing user by ID
func UpdateUserGORM(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := db.Model(&User{}).Where("id = ?", id).Updates(user).Error; err != nil {
		http.Error(w, "User with this name already exists", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Delete a user
func DeleteUserGORM(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if err := db.Delete(&User{}, id).Error; err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Get all tasks
func GetTasksGORM(w http.ResponseWriter, r *http.Request) {
	var tasks []Task
	if err := db.Preload("Users").Find(&tasks).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(tasks)
}

// Create a new task
func CreateTaskGORM(w http.ResponseWriter, r *http.Request) {
	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := db.Create(&task).Error; err != nil {
		http.Error(w, "Failed to create task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

// Update an existing task by ID
func UpdateTaskGORM(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var task Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Find the task by ID
	existingTask := Task{}
	if err := db.First(&existingTask, id).Error; err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	// Update task fields
	existingTask.Title = task.Title
	existingTask.Description = task.Description
	existingTask.Status = task.Status

	// Update user assignments
	if len(task.Users) > 0 {
		db.Model(&existingTask).Association("Users").Clear()
		for _, user := range task.Users {
			db.Model(&existingTask).Association("Users").Append(&user)
		}
	}

	if err := db.Save(&existingTask).Error; err != nil {
		http.Error(w, "Failed to update task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(existingTask)
}

// Delete a task
func DeleteTaskGORM(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var existingTask Task
	if err := db.First(&existingTask, id).Error; err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	if err := db.Model(&existingTask).Association("Users").Clear(); err != nil {
		http.Error(w, "Failed to clear user associations", http.StatusInternalServerError)
		return
	}

	if err := db.Delete(&existingTask).Error; err != nil {
		http.Error(w, "Failed to delete task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func main() {
	initDB()

	r := mux.NewRouter()

	r.HandleFunc("/users", GetUsersGORM).Methods("GET")
	r.HandleFunc("/users", CreateUserGORM).Methods("POST")
	r.HandleFunc("/users/{id}", UpdateUserGORM).Methods("PUT")
	r.HandleFunc("/users/{id}", DeleteUserGORM).Methods("DELETE")

	r.HandleFunc("/tasks", GetTasksGORM).Methods("GET")
	r.HandleFunc("/tasks", CreateTaskGORM).Methods("POST")
	r.HandleFunc("/tasks/{id}", UpdateTaskGORM).Methods("PUT")
	r.HandleFunc("/tasks/{id}", DeleteTaskGORM).Methods("DELETE")

	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
	}).Handler(r)

	log.Fatal(http.ListenAndServe(":8080", handler))
}
