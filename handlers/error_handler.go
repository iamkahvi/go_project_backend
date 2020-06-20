package handlers

import (
	"log"

	"github.com/gin-gonic/gin"
)

// HandleError : return appropriate error code and abort request
func HandleError(status int, err error, c *gin.Context) {
	log.Println(err)
	switch status {
	case 400:
		log.Println(err)
		c.JSON(400, gin.H{"error": err.Error()})
	case 401:
		log.Println("Unauthorized")
		c.JSON(401, gin.H{"error": "Unauthorized"})
	case 500:
		log.Println("Internal Error")
		c.JSON(500, gin.H{"error": err.Error()})
	}
	c.Abort()
}
