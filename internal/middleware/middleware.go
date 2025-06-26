package middleware

import (
	"errors"
	"log"
	"net/http"
	"todoapp/internal/store"
	"todoapp/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserMiddleware struct {
	UserStore  store.UserStore
	TokenStore store.TokenStore
	Logger     *log.Logger
}

// TODO: make middleware to give the handlerfunctions in the group access to the logged in user, admin? user, and tokens

func SetUser(user *store.User, c *gin.Context) {
	c.Set("user", user)
}

func GetUser(c *gin.Context) *store.User {
	user := &store.User{}
	user, ok := c.Keys["user"].(*store.User)
	if !ok {
		c.Abort()
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "missing user in header"})
		return nil
	}
	return user
}

func (um *UserMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		request := struct {
			Csrf_token string `json:"csrf_token" form:"csrf_token" binding:"required"`
		}{}
		err := c.ShouldBind(&request)
		if err != nil {
			SetUser(store.AnonymousUser, c)
			return
		}
		session_token, err := c.Cookie("session_token")
		if err != nil || session_token == "" || request.Csrf_token == "" {
			SetUser(store.AnonymousUser, c)
			return
		}
		token, err := um.TokenStore.GetToken(utils.HashToken(session_token), utils.HashToken(request.Csrf_token))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				SetUser(store.AnonymousUser, c)
				return
			}
			c.Abort()
			um.Logger.Printf("ERROR: userMiddlewareGetToken: %v\n", err)
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		user, err := um.UserStore.GetUserByID(token.UserID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				SetUser(store.AnonymousUser, c)
				return
			}
			c.Abort()
			um.Logger.Printf("ERROR: userMiddlewareGetUserByID: %v\n", err)
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		SetUser(user, c)
		c.Next()
	}
}

func (um *UserMiddleware) RequreLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := GetUser(c)
		if user.IsAnonymous() {
			c.Abort()
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid session"})
			return
		}
		c.Next()
	}
}
