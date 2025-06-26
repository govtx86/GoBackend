package routes

import (
	"net/http"
	"todoapp/internal/app"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(app *app.Application) http.Handler {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Add your frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true, // Enable cookies/auth
	}))

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World!")
	})

	r.GET("/messages", app.MessageHandler.HandleGetMessage)
	r.POST("/messages", app.MessageHandler.HandleCreateMessage)

	r.POST("/register", app.UserHandler.HandleRegister)
	r.POST("/login", app.UserHandler.HandleLogin)

	{
		auth := r.Group("/")
		auth.Use(app.Middleware.Authenticate())
		auth.POST("/logout", app.UserHandler.HandleLogout)
		{
			reqlogin := auth.Group("/")
			reqlogin.Use(app.Middleware.RequreLogin())
			reqlogin.GET("/protected", app.UserHandler.HandleProtected)
		}
	}
	return r
}
