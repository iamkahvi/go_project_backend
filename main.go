package main

import (
	"fmt"
	"strings"

	"example.com/gin_server/storage"
	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
)

// ResponseBody : struct to bind the JSON response body
type ResponseBody struct {
	User string `json:"user"`
}

const m string = "pong"

func pong(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(200, gin.H{
		"message": "pong",
	})
	fmt.Println(c.Request.Header)
	fmt.Println(c.Request.Body)
}

func showUser(d *storage.DBS) gin.HandlerFunc {
	return func(c *gin.Context) {

		user := c.Params.ByName("name")
		d.Users = append(d.Users, user)

		d.Number++
		// fmt.Println(c.Request.Proto)

		value, ok := d.DB[user]
		if ok {
			c.JSON(200, gin.H{"user": user, "value": value, "something": d.PrintAll()})
		} else {
			c.JSON(200, gin.H{"user": user, "status": "no value"})
		}
	}
}

func changeUserValueString(d *storage.DBS) gin.HandlerFunc {
	return func(c *gin.Context) {
		nameValue := c.Params.ByName("nameValue")
		d.Number++

		// fmt.Println(c.Request)
		params := strings.Split(nameValue, "=")
		user, value := params[0], params[1]

		d.DB[user] = value
		c.JSON(200, gin.H{"message": "saved"})
	}
}

func changeUserValueQuery(d *storage.DBS) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.Query("user")
		value := c.Query("value")
		d.Number++

		d.DB[user] = value
		c.JSON(200, gin.H{"message": "saved"})
	}
}

func post(d *storage.DBS) gin.HandlerFunc {
	return func(c *gin.Context) {
		d.Number++

		var rb ResponseBody
		c.BindJSON(&rb)

		d.Users = append(d.Users, rb.User)
		fmt.Println(rb)

		c.Header("Access-Control-Allow-Origin", "*")
		c.JSON(200, gin.H{"list": d.UserList()})
	}
}

func fetchUserList(d *storage.DBS) gin.HandlerFunc {
	return func(c *gin.Context) {
		d.Number++
		c.Header("Access-Control-Allow-Origin", "*")
		c.JSON(200, gin.H{"list": d.UserList()})
	}
}

func main() {
	d := storage.DBS{Number: 0, DB: make(map[string]string)}

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
	r.GET("/users", fetchUserList(&d))

	r.GET("/user/:name", showUser(&d))

	r.GET("/change/:nameValue", changeUserValueString(&d))

	r.GET("/params", changeUserValueQuery(&d))

	r.POST("/post", post(&d))

	r.Run(":8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
