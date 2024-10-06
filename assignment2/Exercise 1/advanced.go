package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "20120117"
	dbname   = "assignment2part2"
)

func createTable(db *sql.DB) {
	tableQuery := `CREATE TABLE IF NOT EXISTS users(
		ID SERIAL PRIMARY KEY,
		name VARCHAR(50) UNIQUE NOT NULL,
		age INT NOT NULL
	);`
	_, err := db.Exec(tableQuery)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Table created")
}

func addUsersWithTransaction(users []map[string]interface{}, db *sql.DB) {
	transaction, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	// shorten transaction for easy use
	ts := transaction

	stringQuery, err := ts.Prepare("INSERT INTO users (name, age) VALUES ($1, $2)")

	if err != nil {
		ts.Rollback()
		log.Fatal(err)
	}
	defer stringQuery.Close()

	for _, user := range users {
		_, err = stringQuery.Exec(user["name"], user["age"])
		if err != nil {
			ts.Rollback()
			log.Fatalf("FAIL", err)
		}
	}
	err = ts.Commit()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Success in adding users")
}

func queryData(ageFilter *int, page, pageSize int, db *sql.DB) {
	var rows *sql.Rows
	var err error
	if ageFilter != nil {
		rows, err = db.Query("SELECT ID, name, age FROM users WHERE age = $1 LIMIT $2 OFFSET $3", *ageFilter, pageSize, (page-1)*pageSize)
	} else {
		rows, err = db.Query("SELECT ID, name, age FROM users LIMIT $1 OFFSET $2", pageSize, (page-1)*pageSize)
	}

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println("Your users:")

	for rows.Next() {
		var ID int
		var name string
		var age int
		if err := rows.Scan(&ID, &name, &age); err != nil {
			log.Fatal(err)
		}
		fmt.Println(ID, name, age)
	}
}

func updateUser(ID int, newName string, newAge int, db *sql.DB) {
	result, err := db.Exec("UPDATE users SET name = $1, age = $2 WHERE id = $3", newName, newAge, ID)
	if err != nil {
		log.Fatal(err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	if affected == 0 {
		fmt.Println("No user found with id:", ID)
	} else {
		fmt.Println("Updated user with id:", ID)
	}
}

func deleteUser(id int, db *sql.DB) {
	result, err := db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		log.Fatalf("Unable to delete user: %v", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		log.Fatalf("Error checking affected rows: %v", err)
	}
	if affected == 0 {
		fmt.Printf("No user with id %d found.\n", id)
	} else {
		fmt.Printf("User with id %d deleted successfully!\n", id)
	}
}

func main() {
	fmt.Println("This is the work of Baatar")

	login := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", login)
	if err != nil {
		log.Fatal("FAIL", err)
	}
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	fmt.Println("Success in connecting")

	// createTable(db)
	// users := []map[string]interface{}{
	// 	{"name": "GEBE", "age": 22},
	// 	{"name": "BAATAR", "age": 22},
	// 	{"name": "GBJE", "age": 22},
	// }
	// addUsersWithTransaction(users, db)

	queryData(nil, 1, 2, db)
	queryData(nil, 2, 2, db)

	ageFilter := 23

	queryData(&ageFilter, 1, 2, db)

	updateUser(1, "Almas", 25, db)

	deleteUser(3, db)

}
