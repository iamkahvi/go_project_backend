package handlers

import (
	"crypto/rand"
	"errors"
	"log"
	"math/big"
	"os"
	"time"

	"example.com/gin_server/email"
	"example.com/gin_server/storage"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// CodeExpiryDuration : how long a passcode is valid for
var CodeExpiryDuration, _ = time.ParseDuration("5m")

// AuthReqBody : struct to bind the JSON response body
type AuthReqBody struct {
	Code  int    `json:"code"`
	Email string `json:"email"`
}

// HandleCodeReq : sends passcode to provided email
func HandleCodeReq(d *storage.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		d.Number++

		to := c.Param("email")

		r, err := generateCode()
		if err != nil {
			log.Println("Generating passcode failed")
			HandleError(500, err, c)
			return
		}

		d.CodeMap[to] = storage.CodeItem{
			Code:   r,
			Expiry: time.Now().Add(CodeExpiryDuration),
		}

		err = email.SendCode(to, r)
		if err != nil {
			log.Println("Sending email failed")
			HandleError(500, err, c)
			return
		}

		c.Header("Access-Control-Allow-Origin", "*")
		c.JSON(200, gin.H{"status": "success"})
	}
}

// HandleCodeSubmit : returns a JWT token if passcode is correct
func HandleCodeSubmit(d *storage.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		d.Number++
		c.Header("Access-Control-Allow-Origin", "*")

		var rb AuthReqBody
		c.BindJSON(&rb)

		code, exp := d.CodeMap[rb.Email].Code, d.CodeMap[rb.Email].Expiry

		if exp.Before(time.Now()) {
			HandleError(400, errors.New("Passcode expired"), c)
			return
		}

		if code < 0 {
			HandleError(500, errors.New("Invalid internal code"), c)
			return
		}

		if code == rb.Code {
			log.Println("Valid Code")
			token, err := generateJWT(rb.Email, []byte(os.Getenv("SECRET")))
			if err != nil {
				log.Println("Error generating JWT")
				HandleError(500, err, c)
				return

			}
			c.JSON(200, gin.H{"status": "success", "token": token})
			return
		}

		HandleError(401, errors.New("Invalid code"), c)
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
