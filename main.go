package main

import (
	"math"
	"net/http"
	"os"
	"strconv"

	"example.com/gin_server/routes"
	"example.com/gin_server/storage"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var db storage.DB

func main() {
	os.Setenv("SECRET", "verysecretkey")
	db.InitDB()

	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowHeaders = []string{"authorization", "content-type"}
	config.AllowOrigins = []string{"http://localhost:3000"}
	r.Use(cors.New(config))

	routes.InitializeRoutes(r, &db)

	r.GET("/ping", pong)

	r.GET("/calc", handleCalculate)

	r.Run()
}

func handleCalculate(c *gin.Context) {
	variableX, _ := strconv.Atoi(c.Query("x"))
	variableY, _ := strconv.Atoi(c.Query("y"))

	answer := calculateSomething(variableX, variableY)

	c.JSON(http.StatusAccepted, gin.H{"answer": answer})
}

func calculateSomething(vx int, vy int) float64 {
	res1 := math.Yn(vx, float64(vy))
	return res1
}

const m string = "pong"

func pong(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
