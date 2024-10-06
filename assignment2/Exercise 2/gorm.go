package main

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "20120117"
	dbname   = "assignment2part2"
)

type User struct {
	ID      int `gorm:"primaryKey"`
	Name    string
	Age     int
	Profile Profile `gorm:"constraint:OnDelete:CASCADE;"`
}

type Profile struct {
	ID                int `gorm:"primaryKey"`
	UserID            int `gorm:"index,unique"`
	Bio               string
	ProfilePictureURL string
}

func insertUserProfile(db *gorm.DB, name string, age int, bio string, profilePictureURL string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		user := User{
			Name: name,
			Age:  age,
		}
		if err := tx.Create(&user).Error; err != nil {
			return err
		}
		Profile := Profile{
			UserID:            user.ID,
			Bio:               bio,
			ProfilePictureURL: profilePictureURL,
		}
		if err := tx.Create(&Profile).Error; err != nil {
			return err
		}
		return nil
	})
}

func getUser(db *gorm.DB) ([]User, error) {
	var userTable []User

	if err := db.Preload("Profile").Find(&userTable).Error; err != nil {
		return nil, err
	}
	return userTable, nil
}

func updateUserProfile(db *gorm.DB, userID int, bio string, URL string) error {
	return db.Model(&Profile{}).Where("user_id = ?", userID).Updates(Profile{
		Bio:               bio,
		ProfilePictureURL: URL,
	}).Error
}

func deleteUser(db *gorm.DB, userID int) error {
	return db.Delete(&User{}, userID).Error
}

func main() {
	fmt.Println("This is the work of Baatar")

	login := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := gorm.Open(postgres.Open(login), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	poolingDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}

	poolingDB.SetMaxIdleConns(5)
	poolingDB.SetMaxOpenConns(10)
	poolingDB.SetConnMaxIdleTime(30 * time.Minute)

	// Testing for error
	if err := db.AutoMigrate(&User{}, &Profile{}); err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}

	// if err := insertUserProfile(db, "gebe", 21, "He is GEBE", "gebe.jpg"); err != nil {
	// 	log.Fatal(err)
	// }
	// if err := insertUserProfile(db, "joe", 35, "He is joe", "joe.jpg"); err != nil {
	// 	log.Fatal(err)
	// }
	// if err := insertUserProfile(db, "jane", 19, "She is jane", "jane.jpg"); err != nil {
	// 	log.Fatal(err)
	// }
	// if err := insertUserProfile(db, "doe", 69, "He is Doe", "doe.jpg"); err != nil {
	// 	log.Fatal(err)
	// }
	allUsers, err := getUser(db)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Users: ", allUsers)

	// if err := updateUserProfile(db, allUsers[0].ID, "This user doesnt like GO", "UpdatedPicture.jpg"); err != nil {
	// 	log.Fatal(err)
	// }

	// if err := deleteUser(db, 3); err != nil {
	// 	log.Fatal(err)
	// } else {
	// 	fmt.Println("Deleted user")
	// }
}
