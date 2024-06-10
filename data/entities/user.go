package entities

import (
	"github.com/google/uuid"
	"messenger/data/common"
	"time"
)

// UserRef is an entity designed to be used as a reference to a user.
type UserRef string

type User struct {
	ID          UserRef `gorm:"primaryKey"`
	Created     time.Time
	Credentials UserCredentials           `gorm:"embedded;embeddedPrefix:credentials_"`
	AuthMethods UserAuthenticationMethods `gorm:"embedded;embeddedPrefix:auth_"`
	Meta        UserMeta                  `gorm:"embedded;embeddedPrefix:meta_"`
}

type UserCredentials struct {
	Email string
	Name  string
}

// UserAuthenticationMethods contains the methods used to authenticate a user.
type UserAuthenticationMethods struct {
	// PasswordHash is the hash of the password used to authenticate the user.
	PasswordHash string
}

type UserMeta struct {
	Categories UserMetaCategories `gorm:"embedded;embeddedPrefix:categories_"`
}

type UserMetaCategories struct {
	Values common.JsonMap[string, float64]
}

// Implementation:

func NewUser() *User {
	return &User{
		ID:      UserRef(uuid.NewString()),
		Created: time.Now(),
		Meta: UserMeta{
			Categories: UserMetaCategories{
				Values: make(common.JsonMap[string, float64]),
			},
		},
	}
}

func NewUserAuthenticationItem() *UserAuthenticationMethods {
	return &UserAuthenticationMethods{}
}
