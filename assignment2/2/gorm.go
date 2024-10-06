package main

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID     int `gorm:"primaryKey"`
	Name   string
	Age    int
	Gender string
}

func addUser(db *gorm.DB, name string, age int, gender string) {
	newUser := User{Name: name, Age: age, Gender: gender}
	adding := db.Create(&newUser)

	if adding.Error != nil {
		log.Fatal(adding.Error)
	}
	fmt.Println("New user is:", newUser.Name, newUser.Age, newUser.Gender)
}

func queryUsers(db *gorm.DB) {
	var users []User
	result := db.Find(&users)
	if result.Error != nil {
		log.Fatal(result.Error)
	}

	for _, user := range users {
		fmt.Printf("User: %s, Age: %d,Gender: %s \n", user.Name, user.Age, user.Gender)
	}
}

func main() {
	fmt.Println("This is the work of Baatar")
	const (
		host     = "localhost"
		port     = 5432
		user     = "postgres"
		password = "20120117"
		dbname   = "postgres"
	)

	login := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := gorm.Open(postgres.Open(login), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&User{})

	addUser(db, "GBJE", 12, "male")

	queryUsers(db)
}
