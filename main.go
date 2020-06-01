package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"strconv"

	"example.com/gin_server/email"
	"example.com/gin_server/storage"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// email.SendCode("iamkahvi@gmail.com", 1234)

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

	r.POST("/auth/:email", handleCodeReq(&db))

	r.POST("/auth", handleAuthReq(&db))

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
		fmt.Println("error:", err)
		return []byte{}, err
	}
	// The slice should now contain random bytes instead of only zeroes.
	return b, nil
}

// AuthReqBody : struct to bind the JSON response body
type AuthReqBody struct {
	Code  int    `json:"code"`
	Email string `json:"email"`
	Token string `json:"token"`
}

func handleAuthReq(d *storage.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		d.Number++
		c.Header("Access-Control-Allow-Origin", "*")

		var rb AuthReqBody
		c.BindJSON(&rb)

		email := d.TokenMap[rb.Token]
		if email != "" {
			// TODO: Implement JWT secret check
			// token, err := jwt.Parse(rb.Token, func(token *jwt.Token) (interface{}, error) {
			// 	// Don't forget to validate the alg is what you expect:
			// 	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			// 		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			// 	}

			// 	// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			// 	return hmacSampleSecret, nil
			// })

			// if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// 	fmt.Println(claims["foo"], claims["nbf"])
			// } else {
			// 	fmt.Println(err)
			// }

			fmt.Println("Signed In ", email)
			c.JSON(200, gin.H{"status": "success", "token": rb.Token})
			return
		}

		r := d.CodeMap[rb.Email]

		if r < 0 {
			fmt.Println("Invalid internal code")
			c.JSON(500, gin.H{"error": "Invalid internal code"})
			return
		}

		if r == rb.Code {
			fmt.Println("Valid Code")
			token, err := generateJWT(rb.Email, d.TokenMap)
			if err != nil {
				fmt.Println("Error generating JWT")
				c.JSON(500, gin.H{"error": "Internal error"})
				return

			}
			c.JSON(200, gin.H{"status": "success", "token": token})
			return
		}

		fmt.Println("Invalid code")
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
			fmt.Println(err)
			log.Fatalf("Generating passcode failed")
		}
		d.CodeMap[to] = r

		err = email.SendCode(to, r)
		if err != nil {
			fmt.Println(err)
			log.Fatalf("Sending email failed")
		}

		c.Header("Access-Control-Allow-Origin", "*")
		c.JSON(200, gin.H{"status": "success"})
	}
}

func generateJWT(email string, tokens map[string]string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
	})

	// Sign and get the complete encoded token as a string using the secret
	r, err := generateRandBytes()
	if err != nil {
		fmt.Println(err)
		log.Fatalf("Generating bytes failed")
		return "", err
	}

	tokenString, err := token.SignedString(r)
	if err != nil {
		fmt.Println(err)
		log.Fatalf("Generating JWT failed")
		return "", err
	}

	tokens[tokenString] = email

	return tokenString, nil
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
