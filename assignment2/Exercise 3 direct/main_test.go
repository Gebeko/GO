package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var testDB *sql.DB

func init() {
	var err error
	login := "host=localhost port=5432 user=postgres password=20120117 dbname=RESTAPI sslmode=disable"
	testDB, err = sql.Open("postgres", login)
	if err != nil {
		panic("failed to connect to test database")
	}

	// Create the users table if it does not exist
	_, err = testDB.Exec("CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, name VARCHAR(100) UNIQUE NOT NULL, age INT NOT NULL)")
	if err != nil {
		panic("failed to create users table")
	}
}

func clearUsers() {
	_, err := testDB.Exec("DELETE FROM users")
	if err != nil {
		panic("failed to clear users table")
	}
}

func TestMain(m *testing.M) {
	// Run tests
	code := m.Run()
	// Cleanup
	testDB.Close()
	// Exit with code
	os.Exit(code)
}

func TestGetUsers(t *testing.T) {
	clearUsers()

	// Prepare a test user
	testUser := User{Name: "GEBE", Age: 22}
	_, err := testDB.Exec("INSERT INTO users (name, age) VALUES ($1, $2)", testUser.Name, testUser.Age)
	if err != nil {
		t.Fatalf("Could not insert test user: %v", err)
	}

	// Set up the router
	router := mux.NewRouter()
	router.HandleFunc("/users", GetUsers).Methods("GET")

	// Act
	req, err := http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Assert
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var users []User
	if err := json.NewDecoder(rr.Body).Decode(&users); err != nil {
		t.Fatalf("Could not decode response: %v", err)
	}

	if len(users) != 1 || users[0].Name != testUser.Name {
		t.Errorf("Unexpected users returned: got %+v want %+v", users, testUser)
	}
}

func TestCreateUser(t *testing.T) {
	clearUsers()

	// Set up the router
	router := mux.NewRouter()
	router.HandleFunc("/users", CreateUser).Methods("POST")

	newUser := User{Name: "Jane Doe", Age: 25}
	userJSON, _ := json.Marshal(newUser)

	// Act
	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(userJSON))
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Assert
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	// Check the user was created in the database
	var createdUser User
	err = testDB.QueryRow("SELECT * FROM users WHERE name = $1", newUser.Name).Scan(&createdUser.ID, &createdUser.Name, &createdUser.Age)
	if err != nil {
		t.Fatalf("Could not retrieve created user: %v", err)
	}

	if createdUser.Name != newUser.Name {
		t.Errorf("Unexpected user created: got %+v want %+v", createdUser, newUser)
	}
}

func TestUpdateUser(t *testing.T) {
	clearUsers()

	// Prepare a test user
	testUser := User{Name: "GEBE", Age: 22}
	result, err := testDB.Exec("INSERT INTO users (name, age) VALUES ($1, $2)", testUser.Name, testUser.Age)
	if err != nil {
		t.Fatalf("Could not insert test user: %v", err)
	}

	// Get the user ID
	userID, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("Could not retrieve last inserted ID: %v", err)
	}

	// Set up the router
	router := mux.NewRouter()
	router.HandleFunc("/users/{id}", UpdateUser).Methods("PUT")

	// Update the user's name
	updatedUser := User{Name: "GEBE BAATAR", Age: 25}
	userJSON, _ := json.Marshal(updatedUser)

	// Act
	req, err := http.NewRequest("PUT", "/users/"+strconv.FormatInt(userID, 10), bytes.NewBuffer(userJSON))
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Assert
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the user was updated in the database
	var resultUser User
	err = testDB.QueryRow("SELECT * FROM users WHERE id = $1", userID).Scan(&resultUser.ID, &resultUser.Name, &resultUser.Age)
	if err != nil {
		t.Fatalf("Could not retrieve updated user: %v", err)
	}

	if resultUser.Name != updatedUser.Name {
		t.Errorf("Unexpected user updated: got %+v want %+v", resultUser, updatedUser)
	}
}

func TestDeleteUser(t *testing.T) {
	clearUsers()

	// Prepare a test user
	testUser := User{Name: "Gebek", Age: 112}
	result, err := testDB.Exec("INSERT INTO users (name, age) VALUES ($1, $2)", testUser.Name, testUser.Age)
	if err != nil {
		t.Fatalf("Could not insert test user: %v", err)
	}

	// Get the user ID
	userID, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("Could not retrieve last inserted ID: %v", err)
	}

	// Set up the router
	router := mux.NewRouter()
	router.HandleFunc("/users/{id}", DeleteUser).Methods("DELETE")

	// Act
	req, err := http.NewRequest("DELETE", "/users/"+strconv.FormatInt(userID, 10), nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Assert
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check that the user was deleted
	var existingUser User
	err = testDB.QueryRow("SELECT * FROM users WHERE id = $1", userID).Scan(&existingUser.ID, &existingUser.Name, &existingUser.Age)
	if err == nil {
		t.Error("User was not deleted")
	}
}

func TestCreateUserDuplicate(t *testing.T) {
	clearUsers()

	// Prepare a test user
	testUser := User{Name: "Baatar", Age: 123}
	_, err := testDB.Exec("INSERT INTO users (name, age) VALUES ($1, $2)", testUser.Name, testUser.Age)
	if err != nil {
		t.Fatalf("Could not insert test user: %v", err)
	}

	// Set up the router
	router := mux.NewRouter()
	router.HandleFunc("/users", CreateUser).Methods("POST")

	duplicateUser := User{Name: "Baatar", Age: 123}
	userJSON, _ := json.Marshal(duplicateUser)

	// Act
	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(userJSON))
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Assert
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}
