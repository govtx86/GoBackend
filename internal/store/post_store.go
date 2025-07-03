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
	GetAllPosts() ([]Post, error)
	GetPostByID(uuid.UUID) (*Post, error)
}

func (pg *PostgresPostStore) CreatePost(post *Post) error {
	result := pg.db.Create(post)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (pg *PostgresPostStore) GetAllPosts() ([]Post, error) {
	var posts []Post
	result := pg.db.Select("id", "title").Find(&posts)
	if result.Error != nil {
		return nil, result.Error
	}
	return posts, nil
}

func (pg *PostgresPostStore) GetPostByID(id uuid.UUID) (*Post, error) {
	var post Post
	result := pg.db.First(&post, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &post, nil
}

