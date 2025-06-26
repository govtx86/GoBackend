package app

import (
	"log"
	"os"
	"todoapp/internal/api"
	"todoapp/internal/middleware"
	"todoapp/internal/store"

	"gorm.io/gorm"
)

type Application struct {
	Logger         *log.Logger
	MessageHandler *api.MessageHandler
	UserHandler    *api.UserHandler
	Middleware     middleware.UserMiddleware
	DB             *gorm.DB
}

func NewApplication() (*Application, error) {
	pgDB, err := store.Open()
	if err != nil {
		return nil, err
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	messageStore := store.NewPostgresMessageStore(pgDB)
	userStore := store.NewPostgresUserStore(pgDB)
	tokenStore := store.NewPostgresTokenStore(pgDB)

	messageHandler := api.NewMessageHandler(messageStore, logger)
	userHandler := api.NewUserHanlder(userStore, tokenStore, logger)

	userMidleware := middleware.UserMiddleware{
		UserStore:  userStore,
		TokenStore: tokenStore,
		Logger:     logger,
	}

	app := &Application{
		Logger:         logger,
		MessageHandler: messageHandler,
		UserHandler:    userHandler,
		Middleware:     userMidleware,
		DB:             pgDB,
	}
	return app, nil
}
