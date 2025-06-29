package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PostHandler struct {
	logger *log.Logger
}

func NewPostHanlder(logger *log.Logger) *PostHandler {
	return &PostHandler{
		logger: logger,
	}
}

func (ph *PostHandler) HandleUploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "image upload error"})
	}
	filename := uuid.New().String()
	c.SaveUploadedFile(file, "./files/static/images/"+filename)
	c.JSON(http.StatusOK, gin.H{"url": fmt.Sprintf("http://%v/static/images/%v", c.Request.Host, filename)})
}
