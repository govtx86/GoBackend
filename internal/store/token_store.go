package store

import (
	"time"
	"todoapp/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Token struct {
	ID           int       `json:"-"`
	UserID       uuid.UUID `gorm:"not null;"`
	User         User      `gorm:"constraint:OnDelete:CASCADE;"`
	SessionToken TokenItem `gorm:"embedded;embeddedPrefix:session_token_" json:"-"`
	CSRFToken    TokenItem `gorm:"embedded;embeddedPrefix:csrf_token_" json:"-"`
	CreatedAt    time.Time `json:"-"`
	UpdatedAt    time.Time `json:"-"`
}

type TokenItem struct {
	PlainText string `gorm:"-" json:"-"`
	Hash      string `json:"-"`
}

type PostgresTokenStore struct {
	db *gorm.DB
}

func NewPostgresTokenStore(db *gorm.DB) *PostgresTokenStore {
	err := db.AutoMigrate(&Token{})
	if err != nil {
		panic(err)
	}
	return &PostgresTokenStore{
		db: db,
	}
}

type TokenStore interface {
	CreateToken(uuid.UUID) (*Token, error)
	GetToken(session_token string, csrf_token string) (*Token, error)
	DeleteAllTokenForUser(uuid.UUID) error
}

func (pg *PostgresTokenStore) CreateToken(userId uuid.UUID) (*Token, error) {
	token := &Token{
		UserID: userId,
	}
	var err error
	token.SessionToken.PlainText, err = utils.GenerateToken(32)
	if err != nil {
		return nil, err
	}
	token.CSRFToken.PlainText, err = utils.GenerateToken(32)
	if err != nil {
		return nil, err
	}
	token.SessionToken.Hash = utils.HashToken(token.SessionToken.PlainText)
	token.CSRFToken.Hash = utils.HashToken(token.CSRFToken.PlainText)
	result := pg.db.Create(token)
	if result.Error != nil {
		return nil, result.Error
	}
	return token, nil
}

func (pg *PostgresTokenStore) GetToken(session_token string, csrf_token string) (*Token, error) {
	token := &Token{}
	result := pg.db.Where("session_token_hash = ? AND csrf_token_hash = ?", session_token, csrf_token).Find(&token)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected != 1 {
		return nil, gorm.ErrRecordNotFound
	}
	return token, nil
}

func (pg *PostgresTokenStore) DeleteAllTokenForUser(userId uuid.UUID) error {
	result := pg.db.Where("user_id = ?", userId).Delete(&Token{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}
