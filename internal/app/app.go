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
	UserHandler    *api.UserHandler
	PostHandler *api.PostHandler
	Middleware     middleware.UserMiddleware
	DB             *gorm.DB
}

func NewApplication() (*Application, error) {
	pgDB, err := store.Open()
	if err != nil {
		return nil, err
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)


	userStore := store.NewPostgresUserStore(pgDB)
	tokenStore := store.NewPostgresTokenStore(pgDB)
	postStore := store.NewPostgresPostStore(pgDB)

	userHandler := api.NewUserHanlder(userStore, tokenStore, logger)
	postHandler := api.NewPostHanlder(postStore, logger)

	userMidleware := middleware.UserMiddleware{
		UserStore:  userStore,
		TokenStore: tokenStore,
		Logger:     logger,
	}

	app := &Application{
		Logger:         logger,
		UserHandler:    userHandler,
		PostHandler: postHandler,
		Middleware:     userMidleware,
		DB:             pgDB,
	}
	return app, nil
}
