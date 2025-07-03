package api

import (
	"fmt"
	"log"
	"net/http"
	"todoapp/internal/middleware"
	"todoapp/internal/store"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PostHandler struct {
	postStore store.PostStore
	logger    *log.Logger
}

func NewPostHanlder(postStore store.PostStore, logger *log.Logger) *PostHandler {
	return &PostHandler{
		postStore: postStore,
		logger:    logger,
	}
}

func (ph *PostHandler) HandleCreatePost(c *gin.Context) {
	postRequest := struct {
		Title   string `json:"title" form:"title" binding:"required"`
		Content string `json:"content" form:"content" binding:"required"`
	}{}
	err := c.ShouldBind(&postRequest)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	user := middleware.GetUser(c)
	post := &store.Post{
		UserID:  user.ID,
		Title:   postRequest.Title,
		Content: postRequest.Content,
	}
	err = ph.postStore.CreatePost(post)
	if err != nil {
		ph.logger.Printf("ERROR: handleCreatePostCreatePost: %v\n", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.String(http.StatusCreated, "successfully created post")
}

func (ph *PostHandler) HandleUploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "image upload error"})
		return
	}
	filename := uuid.New().String()
	c.SaveUploadedFile(file, "./files/static/images/"+filename)
	c.JSON(http.StatusOK, gin.H{"url": fmt.Sprintf("http://%v/static/images/%v", c.Request.Host, filename)})
}

func (ph *PostHandler) HandleGetAllPosts(c *gin.Context) {
	posts, err := ph.postStore.GetAllPosts()
	if err != nil {
		ph.logger.Printf("ERROR: handleGetAllPosts: %v\n", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.IndentedJSON(http.StatusOK, posts)
}

func (ph *PostHandler) HandleGetPostByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	post, err := ph.postStore.GetPostByID(id)
	if err != nil {
		ph.logger.Printf("ERROR: handleGetPostByID: %v\n", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.IndentedJSON(http.StatusOK, post)
}