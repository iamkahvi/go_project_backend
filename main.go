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

	// r.GET("/user/:name", showUser(&db))

	// r.GET("/change/:nameValue", changeUserValueString(&db))

	// r.GET("/params", changeUserValueQuery(&db))

	r.POST("/post", post(&db))

	r.Run(":8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

const m string = "pong"

func pong(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(200, gin.H{
		"message": "pong",
	})
	// fmt.Println(c.Request.Header)
	// fmt.Println(c.Request.Body)
}

// func showUser(du *storage.DB) gin.HandlerFunc {
// 	return func(c *gin.Context) {

// 		user := c.Params.ByName("name")
// 		d.Users = append(d.Users, user)

// 		d.Number++
// 		// fmt.Println(c.Request.Proto)

// 		value, ok := d.DB[user]
// 		if ok {
// 			c.JSON(200, gin.H{"user": user, "value": value, "something": d.PrintAll()})
// 		} else {
// 			c.JSON(200, gin.H{"user": user, "status": "no value"})
// 		}
// 	}
// }

// func changeUserValueString(d *storage.DBS) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		nameValue := c.Params.ByName("nameValue")
// 		d.Number++

// 		// fmt.Println(c.Request)
// 		params := strings.Split(nameValue, "=")
// 		user, value := params[0], params[1]

// 		d.DB[user] = value
// 		c.JSON(200, gin.H{"message": "saved"})
// 	}
// }

// func changeUserValueQuery(d *storage.DBS) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		user := c.Query("user")
// 		value := c.Query("value")
// 		d.Number++

// 		d.DB[user] = value
// 		c.JSON(200, gin.H{"message": "saved"})
// 	}
// }

// ResponseBody : struct to bind the JSON response body
type ResponseBody struct {
	User string `json:"user"`
}

func post(d *storage.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var rb ResponseBody
		c.BindJSON(&rb)

		d.AddUser(rb.User)
		fmt.Println(rb)

		c.Header("Access-Control-Allow-Origin", "*")
		c.JSON(200, gin.H{"list": d.GetUsers()})
	}
}

func fetchUserList(d *storage.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		d.Number++
		c.Header("Access-Control-Allow-Origin", "*")
		c.JSON(200, gin.H{"list": d.GetUsers()})
	}
}
