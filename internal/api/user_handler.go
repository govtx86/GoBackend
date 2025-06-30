package api

import (
	"errors"
	"log"
	"net/http"
	"os"
	"regexp"
	"todoapp/internal/middleware"
	"todoapp/internal/store"
	"todoapp/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	userStore  store.UserStore
	tokenStore store.TokenStore
	logger     *log.Logger
}

func NewUserHanlder(userStore store.UserStore, tokenStore store.TokenStore, logger *log.Logger) *UserHandler {
	return &UserHandler{
		userStore:  userStore,
		tokenStore: tokenStore,
		logger:     logger,
	}
}

type RegisterUserRequest struct {
	Username string `json:"username" form:"username" binding:"required"`
	Email    string `json:"email" form:"email" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

var frontendURL = os.Getenv("FRONTEND_URL")

func (uh *UserHandler) HandleRegister(c *gin.Context) {
	request := RegisterUserRequest{}
	err := c.ShouldBind(&request)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	err = validateRegisterRequest(request)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	check, err := uh.userStore.DoesUsernameExist(request.Username)
	if err != nil {
		uh.logger.Printf("ERROR: doesUsernameExist: %v\n", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	if check {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "username already in use"})
		return
	}

	hashedPassword, err := utils.HashPassword(request.Password)

	if err != nil {
		uh.logger.Printf("ERROR: registerUserHashPassword: %v\n", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	user := &store.User{
		Username:     request.Username,
		Email:        request.Email,
		PasswordHash: hashedPassword,
	}
	err = uh.userStore.CreateUser(user)
	if err != nil {
		uh.logger.Printf("ERROR: createUser: %v\n", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "internal server errror"})
		return
	}
	c.String(http.StatusCreated, "Successfully created user!")

}

func validateRegisterRequest(request RegisterUserRequest) error {
	if request.Username == "" {
		return errors.New("username is required")
	}
	if len(request.Username) > 50 {
		return errors.New("username cannot be greater than 50 characters")
	}
	if len(request.Username) < 5 {
		return errors.New("username cannot be smaller than 5 characters")
	}
	if request.Email == "" {
		return errors.New("email is required")
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(request.Email) {
		return errors.New("invalid email format")
	}
	if request.Password == "" {
		return errors.New("password is required")
	}
	return nil
}

func (uh *UserHandler) HandleLogin(c *gin.Context) {
	request := struct {
		Username string `json:"username" form:"username" binding:"required"`
		Password string `json:"password" form:"password" binding:"required"`
	}{}
	err := c.ShouldBind(&request)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	user, err := uh.userStore.GetUserByUsername(request.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "user not registerd"})
			return
		}
		uh.logger.Printf("ERROR: loginGetUserByUsername: %v\n", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	if !utils.CheckPasswordHash(request.Password, user.PasswordHash) {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "invalid password"})
		return
	}

	tokens, err := uh.tokenStore.CreateToken(user.ID)
	if err != nil {
		uh.logger.Printf("ERROR: createToken: %v\n", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.SetCookie("session_token", tokens.SessionToken.PlainText, 3600, "/", frontendURL, false, true)
	c.SetCookie("csrf_token", tokens.CSRFToken.PlainText, 3600, "/", frontendURL, false, false)
	c.String(http.StatusOK, "Login successful!")

}

func (uh *UserHandler) HandleLogout(c *gin.Context) {
	user := middleware.GetUser(c)
	if user.IsAnonymous() {
		c.String(http.StatusOK, "Already not logged in!")
	} else {
		err := uh.tokenStore.DeleteAllTokenForUser(user.ID)
		if err != nil {
			uh.logger.Printf("ERROR: handleLogoutDeleteToken: %v\n", err)
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		c.SetCookie("session_token", "", -1, "/", frontendURL, false, true)
		c.SetCookie("csrf_token", "", -1, "/", frontendURL, false, false)
		c.String(http.StatusOK, "Logout successful!")
	}
}

func (uh *UserHandler) HandleProtected(c *gin.Context) {
	c.String(http.StatusOK, "Protected Message")
}

func (uh *UserHandler) HandleGetuser(c *gin.Context) {
	user := middleware.GetUser(c)
	c.IndentedJSON(http.StatusOK, gin.H{
		"username": user.Username,
		"email":    user.Email,
	})
}
