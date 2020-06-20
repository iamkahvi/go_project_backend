package handlers

import (
	"errors"
	"strconv"

	"example.com/gin_server/storage"
	"github.com/gin-gonic/gin"
)

// FetchUser : fetches user from db
func FetchUser(d *storage.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		d.Number++
		id := c.Param("id")

		val, err := strconv.Atoi(id)
		if err != nil {
			HandleError(400, err, c)
			return
		}

		user, err := d.GetUser(val)

		if err != nil {
			HandleError(500, err, c)
			return
		}

		c.Header("Access-Control-Allow-Origin", "*")
		c.JSON(200, gin.H{"status": "success", "user": user})
	}
}

// DeletePostBody : struct to bind the JSON response body
type DeletePostBody struct {
	IDs []uint `json:"ids"`
}

// DeleteUser : delete user from db
func DeleteUser(d *storage.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		d.Number++
		id := c.Param("id")

		val, err := strconv.Atoi(id)
		if err != nil {
			HandleError(400, err, c)
			return
		}

		err = d.DeleteUser(uint(val))
		if err != nil {
			HandleError(500, err, c)
			return
		}

		c.Header("Access-Control-Allow-Origin", "*")
		c.JSON(200, gin.H{"status": "success"})
	}
}

// AddPostBody : struct to bind the JSON response body
type AddPostBody struct {
	User string `json:"user"`
	Age  int    `json:"age"`
}

// Add user to db
func AddUser(d *storage.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		d.Number++

		var rb AddPostBody
		c.BindJSON(&rb)

		if rb.User == "" {
			HandleError(400, errors.New("Invalid format"), c)
			return
		}

		err := d.AddUser(rb.User, rb.Age)
		if err != nil {
			HandleError(500, err, c)
			return
		}

		FetchUserList(d)(c)
	}
}

func contains(s []uint, x uint) bool {
	for _, item := range s {
		if item == x {
			return true
		}
	}
	return false
}

// DeleteUserList : delete users from db
func DeleteUserList(d *storage.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		d.Number++

		var rb DeletePostBody
		c.BindJSON(&rb)

		if len(rb.IDs) == 0 || contains(rb.IDs, 0) {
			HandleError(400, errors.New("Invalid format"), c)
			return
		}

		for _, id := range rb.IDs {
			err := d.DeleteUser(id)
			if err != nil {
				HandleError(500, err, c)
				return
			}
		}

		FetchUserList(d)(c)
	}
}

// FetchUserList : fetch users from db
func FetchUserList(d *storage.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		d.Number++

		users, err := d.GetAllUsers()
		if err != nil {
			HandleError(500, err, c)
			return
		}

		c.Header("Access-Control-Allow-Origin", "*")
		c.JSON(200, gin.H{"status": "success", "users": users})
	}
}
