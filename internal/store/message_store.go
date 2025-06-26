package store

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Message struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()" json:"id"`
	Msg       string    `json:"msg" form:"msg" binding:"required"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type PostgresMessageStore struct {
	db *gorm.DB
}

func NewPostgresMessageStore(db *gorm.DB) *PostgresMessageStore {
	err := db.AutoMigrate(&Message{})
	if err != nil {
		panic(err)
	}
	return &PostgresMessageStore{
		db: db,
	}
}

type MessageStore interface {
	GetMessages() (*[]Message, error)
	CreateMessage(msg string) (Message, error)
}

var ErrNoMessages error = errors.New("no messages found")

func (pg *PostgresMessageStore) GetMessages() (*[]Message, error) {
	var messages []Message

	result := pg.db.Find(&messages)

	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected < 1 {
		return nil, ErrNoMessages
	}

	return &messages, nil
}

func (pg *PostgresMessageStore) CreateMessage(msg string) (Message, error) {
	message := Message{
		Msg: msg,
	}
	result := pg.db.Clauses(clause.Returning{}).Create(&message)
	if result.Error != nil {
		return Message{}, result.Error
	}
	return message, nil
}
