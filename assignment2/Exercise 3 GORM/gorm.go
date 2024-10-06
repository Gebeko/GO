package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

type User struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"unique"`
	Age  int    `json:"age"`
}

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "20120117"
	dbname   = "assignment2part2"
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

	db.AutoMigrate(&User{})
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

func main() {
	initDB()

	r := mux.NewRouter()
	r.HandleFunc("/users", GetUsersGORM).Methods("GET")
	r.HandleFunc("/users", CreateUserGORM).Methods("POST")
	r.HandleFunc("/users/{id}", UpdateUserGORM).Methods("PUT")
	r.HandleFunc("/users/{id}", DeleteUserGORM).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", r))
}
