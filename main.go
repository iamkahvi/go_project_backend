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

	r.POST("/add", addHandler(&db))

	r.POST("/delete", deleteHandler(&db))

	r.Run(":8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
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

		d.AddUser(rb.User)
		fmt.Println(rb)

		c.Header("Access-Control-Allow-Origin", "*")
		c.JSON(200, gin.H{"list": d.GetUsers()})
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

		d.DeleteUser(rb.ID)
		fmt.Println(rb)

		c.Header("Access-Control-Allow-Origin", "*")
		c.JSON(200, gin.H{"list": d.GetUsers()})
	}
}

func fetchUserList(d *storage.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		d.Number++

		c.Header("Access-Control-Allow-Origin", "*")
		c.JSON(200, d.GetUsers())
	}
}
