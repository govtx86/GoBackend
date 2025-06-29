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
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true, // Enable cookies/auth
	}))

	r.GET("/status", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World!")
	})

	r.POST("/register", app.UserHandler.HandleRegister)
	r.POST("/login", app.UserHandler.HandleLogin)
	r.MaxMultipartMemory = 8 << 20 // 8 MiB
	r.Static("/static/images", "./files/static/images/")
	r.POST("/images/upload", app.PostHandler.HandleUploadImage)
	{
		auth := r.Group("/")
		auth.Use(app.Middleware.Authenticate())
		auth.POST("/logout", app.UserHandler.HandleLogout)
		{
			reqlogin := auth.Group("/")
			reqlogin.Use(app.Middleware.RequreLogin())
			reqlogin.GET("/user", app.UserHandler.HandleGetuser)
			reqlogin.GET("/protected", app.UserHandler.HandleProtected)
		}
	}
	return r
}
