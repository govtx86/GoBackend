package routes

import (
	"net/http"
	"os"
	"todoapp/internal/app"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(app *app.Application) http.Handler {
	r := gin.Default()
	frontendURL := os.Getenv("FRONTEND_URL")
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{frontendURL}, // Add your frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true, // Enable cookies/auth
	}))

	r.GET("/status", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World!")
	})

	r.MaxMultipartMemory = 8 << 20 // 8 MiB
	r.Static("/static/images", "./files/static/images/")

	r.POST("/register", app.UserHandler.HandleRegister)
	r.POST("/login", app.UserHandler.HandleLogin)
	{
		auth := r.Group("/")
		auth.Use(app.Middleware.Authenticate())
		auth.POST("/logout", app.UserHandler.HandleLogout)
		{
			reqlogin := auth.Group("/")
			reqlogin.Use(app.Middleware.RequreLogin())
			reqlogin.GET("/user", app.UserHandler.HandleGetuser)
			reqlogin.GET("/protected", app.UserHandler.HandleProtected)

			reqlogin.POST("/posts/image/upload", app.PostHandler.HandleUploadImage)
			reqlogin.POST("/posts/new", app.PostHandler.HandleCreatePost)
		}
	}
	
	return r
}
