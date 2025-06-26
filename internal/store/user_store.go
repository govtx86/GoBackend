package store

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();" json:"id"`
	Username     string    `gorm:"unique; not null;" json:"username"`
	Email        string    `gorm:"not null;" json:"email"`
	PasswordHash string    `gorm:"not null;type:varchar(255)" json:"-"`
	CreatedAt    time.Time `json:"-"`
	UpdatedAt    time.Time `json:"-"`
}

var AnonymousUser = &User{}

func (user *User) IsAnonymous() bool {
	return user == AnonymousUser
}

type PostgresUserStore struct {
	db *gorm.DB
}

func NewPostgresUserStore(db *gorm.DB) *PostgresUserStore {
	err := db.AutoMigrate(&User{})
	if err != nil {
		panic(err)
	}
	return &PostgresUserStore{
		db: db,
	}
}

type UserStore interface {
	DoesUsernameExist(string) (bool, error)
	GetUserByUsername(string) (*User, error)
	GetUserByID(uuid.UUID) (*User, error)
	CreateUser(*User) error
}

func (pg *PostgresUserStore) DoesUsernameExist(username string) (bool, error) {
	user := User{}
	result := pg.db.Where("username = ?", username).Find(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, result.Error
	}
	if user.Username != "" {
		return true, nil
	} else {
		return false, nil
	}
}

func (pg *PostgresUserStore) GetUserByUsername(username string) (*User, error) {
	user := &User{}
	result := pg.db.Where("username = ?", username).Find(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected != 1 {
		return nil, gorm.ErrRecordNotFound
	}
	return user, nil
}

func (pg *PostgresUserStore) GetUserByID(userId uuid.UUID) (*User, error) {
	user := &User{}
	result := pg.db.Where("id = ?", userId).Find(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected != 1 {
		return nil, gorm.ErrRecordNotFound
	}
	return user, nil
}

func (pg *PostgresUserStore) CreateUser(user *User) error {
	result := pg.db.Create(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
