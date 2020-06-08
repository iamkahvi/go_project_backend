package storage

import (
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// Model : definition of gorm model
type Model struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// Item : a struct to store a key and value combo
type Item struct {
	key string
}

// User : The user object for storing a user's info
type User struct {
	gorm.Model
	Name    string
	Age     int
	Message string
}

// CodeItem : The object for storing passcode with expiry date
type CodeItem struct {
	Code   int
	Expiry time.Time
}

// DB : The database struct to store stuff from the server
type DB struct {
	User    User
	Item    Item
	Number  int
	db      *gorm.DB
	CodeMap map[string]CodeItem
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
	d.CodeMap = make(map[string]CodeItem)
	d.CodeMap[""] = CodeItem{Code: -1}
}

// DeleteUser : Method to delete a user
func (d *DB) DeleteUser(id uint) error {
	var user User
	user.Model.ID = id
	err := d.db.Delete(&user).Error

	if err != nil {
		return err
	}

	return nil
}

// AddUser : Method to add user to db
func (d *DB) AddUser(name string, age int) error {
	err := d.db.Create(&User{Name: name, Age: age}).Error
	if err != nil {
		return err
	}

	return nil
}

// PrintAll : Method to output the number field
func (d *DB) PrintAll() map[string]interface{} {
	return gin.H{
		"number": 10,
	}
}

// GetAllUsers : Method to get list of users
func (d *DB) GetAllUsers() ([]User, error) {
	var users []User
	err := d.db.Find(&users).Error
	if err != nil {
		return users, err
	}

	return users, nil
}

// GetUser : Method to get user from id
func (d *DB) GetUser(id int) (User, error) {
	var user User
	err := d.db.Find(&user, id).Error
	if err != nil {
		return user, err
	}

	return user, nil
}
