package storage

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// Item : a struct to store a key and value combo
type Item struct {
	key string
}

// User : The user object for storing a user's info
type User struct {
	gorm.Model
	Name string
}

// DB : The database struct to store stuff from the server
type DB struct {
	User   User
	Item   Item
	Number int
	db     *gorm.DB
}

// InitDB : This gets the sql thing idk
func (d *DB) InitDB() {
	db, err := gorm.Open("mysql", "root@/users?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect to the database :(")
	}
	// defer db.Close()

	db.AutoMigrate(&User{})

	d.db = db
}

// DeleteUser : Method to delete a user
func (d *DB) DeleteUser(id uint) int {
	var user User
	user.Model.ID = id
	d.db.Delete(&user)
	return 1
}

// AddUser : Method to add user to db
func (d *DB) AddUser(name string) int {
	d.db.Create(&User{Name: name})
	var user User
	d.db.First(&user, "name = ?", "Kahvi")

	return 1
}

// PrintAll : Method to output the number field
func (d *DB) PrintAll() map[string]interface{} {
	return gin.H{
		"number": 10,
	}
}

// GetUsers : Method to get list of users
func (d *DB) GetUsers() map[string]interface{} {
	var users []User
	d.db.Find(&users)

	return gin.H{
		"users": users,
	}
}
