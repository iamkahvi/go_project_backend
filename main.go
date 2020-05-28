package main

import (
	"fmt"

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

	r.POST("/users", addHandler(&db))

	r.DELETE("/users", deleteHandler(&db))

	r.Run()
}

const m string = "pong"

func pong(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

// AddPostBody : struct to bind the JSON response body
type AddPostBody struct {
	User string `json:"user"`
}

func addHandler(d *storage.DB) gin.HandlerFunc {
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

		fmt.Println(rb)

		users, err := d.GetUsers()
		if err != nil {
			fmt.Println(err)
			c.Header("Access-Control-Allow-Origin", "*")
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.Header("Access-Control-Allow-Origin", "*")
		c.JSON(200, gin.H{"list": users})
	}
}

// DeletePostBody : struct to bind the JSON response body
type DeletePostBody struct {
	ID uint `json:"id"`
}

func deleteHandler(d *storage.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		d.Number++

		var rb DeletePostBody
		c.BindJSON(&rb)

		err := d.DeleteUser(rb.ID)
		if err != nil {
			fmt.Println(err)
			c.Header("Access-Control-Allow-Origin", "*")
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		fmt.Println(rb)

		users, err := d.GetUsers()
		if err != nil {
			fmt.Println(err)
			c.Header("Access-Control-Allow-Origin", "*")
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.Header("Access-Control-Allow-Origin", "*")
		c.JSON(200, gin.H{"list": users})
	}
}

func fetchUserList(d *storage.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		d.Number++

		users, err := d.GetUsers()
		if err != nil {
			fmt.Println(err)
			c.Header("Access-Control-Allow-Origin", "*")
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.Header("Access-Control-Allow-Origin", "*")
		c.JSON(200, users)
	}
}
