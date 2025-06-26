package api

import (
	"errors"
	"log"
	"net/http"
	"todoapp/internal/store"

	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	messageStore store.MessageStore
	logger       *log.Logger
}

func NewMessageHandler(messageStore store.MessageStore, logger *log.Logger) *MessageHandler {
	return &MessageHandler{
		messageStore: messageStore,
		logger:       logger,
	}
}

func (mh *MessageHandler) HandleGetMessage(c *gin.Context) {
	messages, err := mh.messageStore.GetMessages()

	if err != nil {
		if errors.Is(err, store.ErrNoMessages) {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "no messages"})
			return
		}
		mh.logger.Printf("ERROR: handleGetMessage: %v\n", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": messages})
}

func (mh *MessageHandler) HandleCreateMessage(c *gin.Context) {
	message := store.Message{}

	err := c.ShouldBind(&message)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid create request"})
		return
	}
	message, err = mh.messageStore.CreateMessage(message.Msg)
	if err != nil {
		mh.logger.Printf("ERROR: handleCreateMessage: %v\n", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": message})
}