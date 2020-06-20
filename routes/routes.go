package routes

import (
	"example.com/gin_server/handlers"
	"example.com/gin_server/middleware"
	"example.com/gin_server/storage"
	"github.com/gin-gonic/gin"
)

// InitializeRoutes : setup all routes
func InitializeRoutes(r *gin.Engine, db *storage.DB) {

	// View for main page
	r.LoadHTMLGlob("views/*")
	r.GET("/", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.HTML(200, "main.html", gin.H{
			"title": "Gin Server",
		})
	})

	r.POST("/auth/:email", handlers.HandleCodeReq(db))

	r.POST("/auth", handlers.HandleCodeSubmit(db))

	auth := r.Group("/", middleware.IsAuthorized(db))
	{
		auth.GET("/users", handlers.FetchUserList(db))

		auth.POST("/users", handlers.AddUser(db))

		auth.DELETE("/users", handlers.DeleteUserList(db))

		auth.GET("/users/:id", handlers.FetchUser(db))

		auth.DELETE("/users/:id", handlers.DeleteUser(db))
	}

}
