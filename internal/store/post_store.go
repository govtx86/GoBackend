package store

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Post struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();" json:"id"`
	UserID    uuid.UUID `gorm:"not null;" json:"-"`
	User      User      `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type PostgresPostStore struct {
	db *gorm.DB
}

func NewPostgresPostStore(db *gorm.DB) *PostgresPostStore {
	err := db.AutoMigrate(&Post{})
	if err != nil {
		panic(err)
	}
	return &PostgresPostStore{
		db: db,
	}
}

type PostStore interface {
	CreatePost(*Post) error
}

func (pg *PostgresPostStore) CreatePost(post *Post) error {
	result := pg.db.Create(post)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
