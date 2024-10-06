package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var db *sql.DB

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "20120117"
	dbname   = "RESTAPI"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// CONNECT TO DB
func initDB() {
	var err error
	login := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err = sql.Open("postgres", login)
	if err != nil {
		log.Fatal("FAIL", err)
	}
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	fmt.Println("Success in connecting")
}

// Get users
func GetUsers(w http.ResponseWriter, r *http.Request) {
	age := r.URL.Query().Get("age")
	sort := r.URL.Query().Get("sort")

	query := "SELECT * FROM users"
	var conditions []string
	var args []interface{}

	if age != "" {
		conditions = append(conditions, "age = $1")
		args = append(args, age)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	if sort == "asc" {
		query += " ORDER BY name ASC"
	} else if sort == "desc" {
		query += " ORDER BY name DESC"
	}

	// Pagination
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")
	if limit != "" && offset != "" {
		query += fmt.Sprintf(" LIMIT %s OFFSET %s", limit, offset)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Age); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	json.NewEncoder(w).Encode(users)
}

// Creating user
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Check if name is unique
	var existingUser User
	err := db.QueryRow("SELECT * FROM users WHERE name = $1", user.Name).Scan(&existingUser.ID)
	if err == nil {
		http.Error(w, "User with this name already exists", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("INSERT INTO users (name, age) VALUES ($1, $2)", user.Name, user.Age)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Update an user by ID with name uniqueness
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Check if name is unique
	var existingUser User
	err := db.QueryRow("SELECT * FROM users WHERE name = $1 AND id != $2", user.Name, id).Scan(&existingUser.ID)
	if err == nil {
		http.Error(w, "User with this name already exists", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("UPDATE users SET name = $1, age = $2 WHERE id = $3", user.Name, user.Age, id)
	if err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Delete a user by ID
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	// Check if user exists
	var existingUser User
	err := db.QueryRow("SELECT * FROM users WHERE id = $1", id).Scan(&existingUser.ID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	_, err = db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func main() {
	initDB()
	defer db.Close()

	r := mux.NewRouter()
	r.HandleFunc("/users", GetUsers).Methods("GET")
	r.HandleFunc("/users", CreateUser).Methods("POST")
	r.HandleFunc("/users/{id}", UpdateUser).Methods("PUT")
	r.HandleFunc("/users/{id}", DeleteUser).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", r))
}
