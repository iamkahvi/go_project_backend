package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"strings"

	"example.com/gin_server/email"
	"example.com/gin_server/storage"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// JWTSigningKey : key to sign JWT tokens with
// This definitely should be an env variable or something
var JWTSigningKey = []byte("verysecretkey")

func main() {

	db := storage.DB{Number: 0}
	db.InitDB()

	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://google.com", "http://localhost:3000"}
	r.Use(cors.New(config))

	auth := r.Group("/", isAuthorized(&db))
	{
		auth.GET("/users", fetchUserList(&db))

		auth.POST("/users", addUser(&db))

		auth.DELETE("/users", deleteUserList(&db))

		auth.GET("/users/:id", fetchUser(&db))

		auth.DELETE("/users/:id", deleteUser(&db))
	}

	r.LoadHTMLGlob("views/*")
	r.GET("/", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.HTML(200, "main.html", gin.H{
			"title": "Main website",
		})
	})

	r.POST("/auth/:email", handleCodeReq(&db))

	r.POST("/auth", handleCodeSubmit(&db))

	r.GET("/ping", pong)

	r.Run()
}

func isAuthorized(d *storage.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Splitting auth header by spaces
		authHeader := strings.Split(c.GetHeader("Authorization"), " ")

		if len(authHeader) != 2 {
			log.Println("Invalid auth format")
			c.JSON(400, gin.H{"error": "Invalid auth format"})
			c.Abort()
			return
		}

		if authHeader[0] == "Bearer" {
			claims := jwt.MapClaims{}
			token, err := jwt.ParseWithClaims(authHeader[1], claims, func(token *jwt.Token) (interface{}, error) {
				// Check signing method
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}
				// I'm not doing anymore validation stuff here. I'm only using one signing key
				return JWTSigningKey, nil
			})

			if token.Valid {
				fmt.Println("claims", claims)
				log.Println("authorized")
				return
			}

			if err != nil {
				log.Println(err.Error())
			}
		}

		log.Println("unauthorized")
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		c.Abort()
		return
	}
}

const m string = "pong"

func pong(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func generateCode() (int, error) {
	max := big.NewInt(10000)

	r, err := rand.Int(rand.Reader, max)
	if err != nil {
		return 0, err
	}

	return int(r.Int64()), nil
}

func generateRandBytes() ([]byte, error) {
	c := 10
	b := make([]byte, c)
	_, err := rand.Read(b)
	if err != nil {
		log.Println("error:", err)
		return []byte{}, err
	}
	// The slice should now contain random bytes instead of only zeroes.
	return b, nil
}

// AuthReqBody : struct to bind the JSON response body
type AuthReqBody struct {
	Code  int    `json:"code"`
	Email string `json:"email"`
}

func handleCodeSubmit(d *storage.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		d.Number++
		c.Header("Access-Control-Allow-Origin", "*")

		var rb AuthReqBody
		c.BindJSON(&rb)

		r := d.CodeMap[rb.Email]

		if r < 0 {
			log.Println("Invalid internal code")
			c.JSON(500, gin.H{"error": "Invalid internal code"})
			return
		}

		if r == rb.Code {
			log.Println("Valid Code")
			token, err := generateJWT(rb.Email, JWTSigningKey)
			if err != nil {
				log.Println("Error generating JWT")
				c.JSON(500, gin.H{"error": "Internal error"})
				return

			}
			c.JSON(200, gin.H{"status": "success", "token": token})
			return
		}

		log.Println("Invalid code")
		c.JSON(400, gin.H{"error": "Invalid code"})
		return
	}
}

func handleCodeReq(d *storage.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		d.Number++

		to := c.Param("email")

		r, err := generateCode()
		if err != nil {
			log.Println(err)
			log.Println("Generating passcode failed")
		}
		d.CodeMap[to] = r

		err = email.SendCode(to, r)
		if err != nil {
			log.Println(err)
			log.Println("Sending email failed")
		}

		c.Header("Access-Control-Allow-Origin", "*")
		c.JSON(200, gin.H{"status": "success"})
	}
}

func generateJWT(email string, key []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
	})

	tokenString, err := token.SignedString(key)

	if err != nil {
		log.Println(err)
		log.Println("Generating JWT failed")
		return "", err
	}

	return tokenString, nil
}

func fetchUser(d *storage.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		d.Number++
		id := c.Param("id")

		val, err := strconv.Atoi(id)
		if err != nil {
			log.Println("Invalid format")
			c.Header("Access-Control-Allow-Origin", "*")
			c.JSON(400, gin.H{"error": "Invalid format"})
			return
		}

		user, err := d.GetUser(val)

		if err != nil {
			log.Println(err)
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
			log.Println("Invalid format")
			c.Header("Access-Control-Allow-Origin", "*")
			c.JSON(400, gin.H{"error": "Invalid format"})
			return
		}

		err = d.DeleteUser(uint(val))
		if err != nil {
			log.Println(err)
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
			log.Println("Invalid format")
			c.Header("Access-Control-Allow-Origin", "*")
			c.JSON(400, gin.H{"error": "Invalid format"})
			return
		}

		err := d.AddUser(rb.User)
		if err != nil {
			log.Println(err)
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
			log.Println("Invalid format")
			c.Header("Access-Control-Allow-Origin", "*")
			c.JSON(400, gin.H{"error": "Invalid format"})
			return
		}

		for _, id := range rb.IDs {
			err := d.DeleteUser(id)
			if err != nil {
				log.Println(err)
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
			log.Println(err)
			c.Header("Access-Control-Allow-Origin", "*")
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.Header("Access-Control-Allow-Origin", "*")
		c.JSON(200, gin.H{"status": "success", "users": users})
	}
}
