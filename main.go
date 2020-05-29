package main

import (
	"fmt"
	"strconv"

	"example.com/gin_server/storage"
	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
)

func main() {
	db := storage.DB{Number: 0}
	db.InitDB()

	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://google.com", "http://localhost:3000"}
	r.Use(cors.New(config))

	r.LoadHTMLGlob("views/*")
	r.GET("/", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.HTML(200, "main.html", gin.H{
			"title": "Main website",
		})
	})

	r.GET("/ping", pong)

	r.GET("/users", fetchUserList(&db))

	r.POST("/users", addUser(&db))

	r.DELETE("/users", deleteUserList(&db))

	r.GET("/users/:id", fetchUser(&db))

	r.DELETE("/users/:id", deleteUser(&db))

	r.Run()
}

const m string = "pong"

func pong(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func fetchUser(d *storage.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		d.Number++
		id := c.Param("id")

		val, err := strconv.Atoi(id)
		if err != nil {
			fmt.Println("Invalid format")
			c.Header("Access-Control-Allow-Origin", "*")
			c.JSON(400, gin.H{"error": "Invalid format"})
			return
		}

		user, err := d.GetUser(val)

		if err != nil {
			fmt.Println(err)
			c.Header("Access-Control-Allow-Origin", "*")
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.Header("Access-Control-Allow-Origin", "*")
		c.JSON(200, gin.H{"status": "success", "user": user})
	}
}

func deleteUser(d *storage.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		d.Number++
		id := c.Param("id")

		val, err := strconv.Atoi(id)
		if err != nil {
			fmt.Println("Invalid format")
			c.Header("Access-Control-Allow-Origin", "*")
			c.JSON(400, gin.H{"error": "Invalid format"})
			return
		}

		err = d.DeleteUser(uint(val))
		if err != nil {
			fmt.Println(err)
			c.Header("Access-Control-Allow-Origin", "*")
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.Header("Access-Control-Allow-Origin", "*")
		c.JSON(200, gin.H{"status": "success"})
	}
}

// AddPostBody : struct to bind the JSON response body
type AddPostBody struct {
	User string `json:"user"`
}

func addUser(d *storage.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		d.Number++

		var rb AddPostBody
		c.BindJSON(&rb)

		if rb.User == "" {
			fmt.Println("Invalid format")
			c.Header("Access-Control-Allow-Origin", "*")
			c.JSON(400, gin.H{"error": "Invalid format"})
			return
		}

		err := d.AddUser(rb.User)
		if err != nil {
			fmt.Println(err)
			c.Header("Access-Control-Allow-Origin", "*")
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		fetchUserList(d)(c)
	}
}

// DeletePostBody : struct to bind the JSON response body
type DeletePostBody struct {
	IDs []uint `json:"ids"`
}

func contains(s []uint, x uint) bool {
	for _, item := range s {
		if item == x {
			return true
		}
	}
	return false
}

func deleteUserList(d *storage.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		d.Number++

		var rb DeletePostBody
		c.BindJSON(&rb)

		if len(rb.IDs) == 0 || contains(rb.IDs, 0) {
			fmt.Println("Invalid format")
			c.Header("Access-Control-Allow-Origin", "*")
			c.JSON(400, gin.H{"error": "Invalid format"})
			return
		}

		for _, id := range rb.IDs {
			err := d.DeleteUser(id)
			if err != nil {
				fmt.Println(err)
				c.Header("Access-Control-Allow-Origin", "*")
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
		}

		fetchUserList(d)(c)
	}
}

func fetchUserList(d *storage.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		d.Number++

		users, err := d.GetAllUsers()
		if err != nil {
			fmt.Println(err)
			c.Header("Access-Control-Allow-Origin", "*")
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.Header("Access-Control-Allow-Origin", "*")
		c.JSON(200, gin.H{"status": "success", "users": users})
	}
}
